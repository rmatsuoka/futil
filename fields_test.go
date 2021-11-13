package futil

import (
	"testing"
)

func TestFields(t *testing.T) {
	tests := []struct {
		input  string
		fields []string
		err    error
	}{
		{``, nil, nil},
		{`read`, []string{`read`}, nil},
		{`  write`, []string{`write`}, nil},
		{`get   -n  foo/bar     abc`, []string{"get", "-n", "foo/bar", "abc"}, nil},
		{`''`, []string{""}, nil},
		{`'walk'`, []string{"walk"}, nil},
		{`glob   '*/*/*.c'`, []string{"glob", "*/*/*.c"}, nil},
		{`hel'lo  worl'd   hello`, []string{"hello  world", "hello"}, nil},
		{`ls   hello\ world`, []string{"ls", "hello world"}, nil},
		{`\'`, []string{`'`}, nil},
		{`'\'`, []string{`\`}, nil},
		{`\\`, []string{`\`}, nil},
		{`rock \'n\' roll`, []string{"rock", "'n'", "roll"}, nil},
		{`read qwerty'asdfgh/world`, nil, errMissingQuote},
		{`grep hello\`, nil, errEscapeEnd},
	}
	for _, test := range tests {
		f, err := Fields([]byte(test.input))
		if err != test.err {
			t.Errorf("%s: got field=%q, err=%q, want field=%q, err=%q",
				test.input, f, err, test.fields, test.err)
			continue
		}
		if len(f) != len(test.fields) {
			t.Errorf("%s: got field=%q, err=%q, want field=%q, err=%q",
				test.input, f, err, test.fields, test.err)
			continue
		}
		for i := range test.fields {
			if f[i] != test.fields[i] {
				t.Errorf("%s: got field=%q, err=%q, want field=%q, err=%q",
					test.input, f, err, test.fields, test.err)
				break
			}
		}
	}
}
