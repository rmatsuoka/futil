package futil

import (
	"fmt"
	"io"
	"io/fs"
)

func globMain(fsys fs.FS, w, ew io.Writer, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing argument")
	}
	m, err := fs.Glob(fsys, args[0])
	if err != nil {
		return err
	}
	for _, p := range m {
		fmt.Fprintln(w, p)
	}
	return nil
}
