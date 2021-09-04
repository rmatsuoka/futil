package gofs

import (
	"fmt"
	"io/fs"
)

func walkMain(fsys fs.FS, args []string) {
	var root string
	if len(args) == 0 {
		root = "."
	} else {
		root = args[0]
	}
	err := fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			warn(err)
		}
		fmt.Println(path)
		return nil
	})
	if err != nil {
		errExit(err)
	}
}
