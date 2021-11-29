# futil

io/fs utils. This package provides a shell-like interface to search fs.FS.

## example
This is an example of `Shell` function in example/dirfu.go

```Go
package main

import (
	"os"

	"github.com/rmatsuoka/futil"
)

func main() {
	if err := futil.Shell(os.DirFS("/"), os.Stdin, os.Stdout, os.Stderr, "zipfu % "); err != nil {
		os.Exit(1)
	}
}
```

```console
$ go build
$ ./dirfu 
dirfu % ls -l
Lrwxrwxrwx           7 Mar  2  2020 bin
drwxr-xr-x        4096 Nov  5 11:56 boot
drwxrwxr-x        4096 Mar  2  2020 cdrom
drwxr-xr-x        5200 Nov 29 15:30 dev
...(omitted)
dirfu % glob usr/bin/*grep
usr/bin/egrep
usr/bin/fgrep
usr/bin/grep
usr/bin/pgrep
dirfu % read usr/bin/egrep
#!/bin/sh
exec grep -E "$@"
dirfu % exit
$  
```


## available commands
```
get [-a] [-n] path [local_path]
	save a file in the local.

glob pattern
	search a filename with a pattern.

ls [-d] [-l] [path]
	like Unix ls command.

read path
	read a file.

walk [root]
	display all files from the root directory.

help
	list available commands.

exit
	exit the shell.
```
