package gofs

import (
	"flag"
	"fmt"
	"io/fs"
	"strconv"
	"time"
)

var (
	lsFlagSet = flag.NewFlagSet("ls", flag.ExitOnError)
	lslFlag   = lsFlagSet.Bool("l", false, "print long format")
	lsdFlag   = lsFlagSet.Bool("d", false, "If arg is a dir, print it instead of its entrys")

	// for long format
	modeWidth = 0
	sizeWidth = 0
)

func lsMain(fsys fs.FS, args []string) {
	if err := lsFlagSet.Parse(args); err != nil {
		errExit(err)
	}

	var p string
	if lsFlagSet.NArg() == 0 {
		p = "."
	} else {
		p = lsFlagSet.Args()[0]
	}

	if err := ls(fsys, p); err != nil {
		errExit(err)
	}
}

func ls(fsys fs.FS, name string) error {
	info, err := fs.Stat(fsys, name)
	if err != nil {
		return err
	}

	var infos []fs.FileInfo
	if !info.IsDir() || *lsdFlag {
		infos = append(infos, info)
	} else {
		dirs, err := fs.ReadDir(fsys, name)
		if err != nil {
			return err
		}
		for _, d := range dirs {
			i, err := d.Info()
			if err != nil {
				warn(err)
				continue
			}
			infos = append(infos, i)
		}
	}

	// future: add sort infos

	if *lslFlag {
		doWidth(infos)
	}

	for _, i := range infos {
		if *lslFlag {
			fmt.Printf("%s\n", lsLongFmt(i))
		} else {
			fmt.Printf("%v\n", i.Name())
		}
	}
	return nil
}

func doWidth(infos []fs.FileInfo) {
	for _, info := range infos {
		m := len(info.Mode().String())
		if modeWidth < m {
			modeWidth = m
		}

		s := len(strconv.FormatInt(info.Size(), 10))
		if sizeWidth < s {
			sizeWidth = s
		}
	}
}

func lsLongFmt(info fs.FileInfo) string {
	t := info.ModTime()
	var timeStr string
	if t.Year() == time.Now().Year() {
		timeStr = t.Format("Jan _2 15:04")
	} else {
		timeStr = t.Format("Jan _2  2006")
	}

	return fmt.Sprintf("%-*s %*d %s %s",
		modeWidth,
		info.Mode().String(),
		sizeWidth,
		info.Size(),
		timeStr,
		info.Name())
}
