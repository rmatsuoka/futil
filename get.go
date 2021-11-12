package futil

import (
	"flag"
	"io"
	"io/fs"
	"os"
	pathpkg "path"
)

func getMain(fsys fs.FS, w, ew io.Writer, args []string) error {
	getFlagSet := flag.NewFlagSet("get", flag.ContinueOnError)
	getFlagSet.SetOutput(ew)
	aFlag := getFlagSet.Bool("a", false, "append the srcfile to the localfile, not overwrite them.")
	nFlag := getFlagSet.Bool("n", false, "do not overwrite an existing file")

	err := getFlagSet.Parse(args)
	if err == flag.ErrHelp {
		return nil
	}
	if err != nil {
		return err
	}

	fargs := getFlagSet.Args()
	if len(fargs) < 1 {
		return errMissArg
	}

	srcName := fargs[0]
	var localName string
	if len(fargs) < 2 {
		localName = pathpkg.Base(srcName)
	} else {
		localName = fargs[1]
	}

	oflag := os.O_WRONLY | os.O_CREATE
	if *aFlag {
		oflag |= os.O_APPEND
	} else {
		oflag |= os.O_TRUNC
	}
	if *nFlag {
		oflag |= os.O_EXCL
	}
	var localFile *os.File
	localFile, err = os.OpenFile(localName, oflag, 0644)
	if err != nil {
		return err
	}

	var srcFile fs.File
	srcFile, err = fsys.Open(srcName)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	_, err = io.Copy(localFile, srcFile)
	if closeErr := localFile.Close(); err == nil {
		err = closeErr
	}

	return err
}
