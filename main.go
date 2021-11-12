package futil

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
)

const (
	CommandList = `get [-a] [-n] path [local_path]
glob pattern
ls [-d] [-l] [path]
read path
walk [root]
help
exit`
)

var (
	// run exit command
	Exit = errors.New("exit")

	errMissArg = errors.New("missing argument")
)

func Eval(fsys fs.FS, w io.Writer, ew io.Writer, args []string) error {
	if len(args) == 0 {
		return errMissArg
	}

	var err error
	switch args[0] {
	case "ls":
		err = lsMain(fsys, w, ew, args[1:])
	case "get":
		err = getMain(fsys, w, ew, args[1:])
	case "glob":
		err = globMain(fsys, w, ew, args[1:])
	case "read":
		err = readMain(fsys, w, ew, args[1:])
	case "walk":
		err = walkMain(fsys, w, ew, args[1:])
	case "help":
		_, err = fmt.Fprintln(ew, CommandList)
	case "exit":
		return Exit
	default:
		err = fmt.Errorf("no such command")
	}
	if err != nil {
		return fmt.Errorf("%s: %v", args[0], err)
	}
	return nil
}
