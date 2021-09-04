package gofs

import (
	"fmt"
	"io/fs"
)

func globMain(fsys fs.FS, args []string) {
	if len(args) == 0 {
		errExit(fmt.Errorf("missing argument"))
	}
	m, err := fs.Glob(fsys, args[0])
	if err != nil {
		errExit(err)
	}
	for _, p := range m {
		fmt.Println(p)
	}
}
