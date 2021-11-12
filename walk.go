package futil

import (
	"fmt"
	"io"
	"io/fs"
)

func walkMain(fsys fs.FS, w io.Writer, ew io.Writer, args []string) error {
	root := "."
	if len(args) > 0 {
		root = args[0]
	}

	return fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(ew, "walk: %v", err)
		}
		_, printErr := fmt.Fprintln(w, path)
		return printErr
	})
}
