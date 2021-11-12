package futil

import (
	"io"
	"io/fs"
)

func readMain(fsys fs.FS, w, ew io.Writer, args []string) error {
	if len(args) < 1 {
		return errMissArg
	}

	f, err := fsys.Open(args[0])
	if err != nil {
		return err
	}
	defer f.Close() // ignore err

	_, err = io.Copy(w, f)
	return err
}
