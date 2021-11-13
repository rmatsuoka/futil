package futil

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"unicode"
	"unicode/utf8"
)

type mustWriter struct {
	w io.Writer
}

func (mw *mustWriter) Write(p []byte) (int, error) {
	n, err := mw.w.Write(p)
	if err != nil {
		panic(err)
	}
	return n, nil
}

// if failed to write error message to ew then panic
func Shell(fsys fs.FS, r io.Reader, w, ew io.Writer, prompt string) error {
	mw := &mustWriter{ew}

	var err error
	s := bufio.NewScanner(r)

	fmt.Fprint(w, prompt)
	for s.Scan() {
		fis, err := Fields(s.Bytes())
		if err != nil {
			fmt.Fprintln(ew, err)
			fmt.Fprint(w, prompt)
			continue
		}
		if len(fis) == 0 {
			fmt.Fprint(w, prompt)
			continue
		}
		err = Eval(fsys, w, mw, fis)
		if err == Exit {
			break
		}
		if err != nil {
			fmt.Fprintln(mw, err)
		}
		fmt.Fprint(w, prompt)
	}
	if sErr := s.Err(); sErr != nil {
		err = sErr
		fmt.Fprintln(mw, err)
	}
	return err
}

var (
	errInvalidUTF8  = errors.New("the encoding UTF-8 is invalid")
	errMissingQuote = errors.New("unmatched '")
	errEscapeEnd    = errors.New("escape character at the end of the string")
)

func Fields(b []byte) ([]string, error) {
	f := newFielder(b)
	var fis []string
	for {
		b, err := f.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		fis = append(fis, string(b))
	}
	return fis, nil
}

type fielder struct {
	b     []byte
	start int
	err   error
}

func newFielder(b []byte) *fielder {
	return &fielder{b: b}
}

func (f *fielder) readQuote() ([]byte, error) {
	start := f.start
	end := f.start
	for f.start < len(f.b) {
		ch, s := utf8.DecodeRune(f.b[f.start:])
		if ch == utf8.RuneError {
			return nil, errInvalidUTF8
		}
		f.start += s

		if ch == '\'' {
			return f.b[start:end], nil
		}
		end += s
	}
	return nil, errMissingQuote
}

func (f *fielder) Next() ([]byte, error) {
	if f.err != nil {
		return nil, f.err
	}

	// skip spaces
	for f.start < len(f.b) {
		ch, s := utf8.DecodeRune(f.b[f.start:])
		if ch == utf8.RuneError {
			f.err = errInvalidUTF8
			return nil, f.err
		}
		if !unicode.IsSpace(ch) {
			break
		}
		f.start += s
	}

	// reach the end
	if f.start >= len(f.b) {
		f.err = io.EOF
		return nil, f.err
	}

	var b []byte
	isEscaped := false
loop:
	for f.start < len(f.b) {
		start := f.start
		ch, s := utf8.DecodeRune(f.b[f.start:])
		if ch == utf8.RuneError {
			f.err = errInvalidUTF8
			return nil, f.err
		}
		f.start += s

		if isEscaped {
			b = append(b, f.b[start:f.start]...)
			isEscaped = false
			continue loop
		}

		if unicode.IsSpace(ch) {
			break loop
		}

		switch {
		case ch == '\\':
			isEscaped = true
		case ch == '\'':
			q, err := f.readQuote()
			if err != nil {
				f.err = err
				return nil, f.err
			}
			b = append(b, q...)
		default:
			b = append(b, f.b[start:f.start]...)
		}
	}
	if isEscaped {
		return nil, errEscapeEnd
	}
	return b, nil
}
