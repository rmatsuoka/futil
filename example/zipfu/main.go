package main

import (
	"archive/zip"
	"fmt"
	"os"

	"github.com/rmatsuoka/futil"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: zipfu zipfile")
		os.Exit(1)
	}
	fsys, err := zip.OpenReader(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "zipfu: %v\n", err)
		os.Exit(1)
	}
	if err := futil.Shell(fsys, os.Stdin, os.Stdout, os.Stderr, "zipfu % "); err != nil {
		os.Exit(1)
	}
}
