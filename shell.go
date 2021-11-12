package futil

import (
	"bufio"
	"strings"
	"fmt"	
	"io"
	"io/fs"
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

	fmt.Print(prompt)
	for s.Scan() {
		err = Eval(fsys, w, mw, strings.Fields(s.Text()))
		if err == Exit {
			break
		}
		if err != nil {
			fmt.Fprintln(mw, err)
		}
		fmt.Print(prompt)
	}
	if sErr := s.Err(); sErr != nil {
		err = sErr
		fmt.Fprintln(mw, err)
	}
	return err
}