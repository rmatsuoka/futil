package main

import (
	"archive/zip"
	"fmt"
	"os"

	"github.com/rmatsuoka/gofs"
)

func main() {
	gofs.UsageArgs = "zipfile cmd [arg...]"

	if len(os.Args) < 2 {
		gofs.Usage()
		os.Exit(1)
	}
	fsys, err := zip.OpenReader(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", gofs.Progname, err)
		os.Exit(1)
	}
	gofs.Main(fsys, os.Args[2:])
}
