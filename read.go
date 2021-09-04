package gofs

import (
	"fmt"
	"io"
	"io/fs"
	"os"
)

func readMain(fsys fs.FS, args []string) {
	if len(args) == 0 {
		errExit(fmt.Errorf("missing argument"))
	}
	for _, fname := range args {
		f, err := fsys.Open(fname)
		if err != nil {
			warn(err)
			exitCode = 1
			continue
		}
		if s, _ := f.Stat(); s.IsDir() {
			warn(fmt.Errorf("%s: is a directory", fname))
			exitCode = 1
			continue
		}
		io.Copy(os.Stdout, f)
		f.Close()
	}
}
