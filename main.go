package gofs

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

var (
	Progname         = filepath.Base(os.Args[0])
	UsageArgs        = "cmd [arg...]"
	Usage     func() = usageExample
	cmdname          = "main"
	exitCode         = 0
)

func usageExample() {
	w := flag.CommandLine.Output()
	fmt.Fprintf(w, "usage: %s %s\n", Progname, UsageArgs)
	fmt.Fprintln(w, `cmds:
	glob 'pattern'
	ls [-l] [-d] [path]
	read path...
	walk [root]
note: pattern at glob shuld be quoted.
`)
}

func Main(fsys fs.FS, args []string) {
	gofsFlag := flag.NewFlagSet(Progname, flag.ExitOnError)
	gofsFlag.Usage = Usage

	if err := gofsFlag.Parse(args); err != nil {
		errExit(err)
	}
	pargs := gofsFlag.Args()
	if len(pargs) == 0 {
		gofsFlag.Usage()
		os.Exit(1)
	}

	cmdname = pargs[0]
	switch pargs[0] {
	case "ls":
		lsMain(fsys, args[1:])
	case "glob":
		globMain(fsys, args[1:])
	case "read":
		readMain(fsys, args[1:])
	case "walk":
		walkMain(fsys, args[1:])
	default:
		errExit(fmt.Errorf("No such command"))
	}
	os.Exit(exitCode)
}

func warn(e error) {
	fmt.Fprintf(os.Stderr, "%s: %s: %v\n", Progname, cmdname, e)
}

func errExit(e error) {
	warn(e)
	if exitCode == 0 {
		exitCode = 1
	}
	os.Exit(exitCode)
}
