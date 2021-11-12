package main

import (
	"os"

	"github.com/rmatsuoka/futil"
)

func main() {
	if err := futil.Shell(os.DirFS("/"), os.Stdin, os.Stdout, os.Stderr, "dirfu % "); err != nil {
		os.Exit(1)
	}
}
