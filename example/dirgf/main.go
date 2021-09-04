package main

import (
	"os"

	"github.com/rmatsuoka/gofs"
)

func main() {
	gofs.Main(os.DirFS("/"), os.Args[1:])
}
