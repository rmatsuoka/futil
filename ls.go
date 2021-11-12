package futil

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"strconv"
	"time"
)

func lsMain(fsys fs.FS, w, ew io.Writer, args []string) error {
	lsFlagSet := flag.NewFlagSet("ls", flag.ContinueOnError)
	lsFlagSet.SetOutput(ew)
	lLsFlag := lsFlagSet.Bool("l", false, "List in long format.")
	dLsFlag := lsFlagSet.Bool("d", false, "If arg is a dir, list it, not its entrys.")

	err := lsFlagSet.Parse(args)
	if err == flag.ErrHelp {
		return nil
	}
	if err != nil {
		return err
	}

	name := "."
	if lsFlagSet.NArg() > 0 {
		name = lsFlagSet.Args()[0]
	}

	info, err := fs.Stat(fsys, name)
	if err != nil {
		return err
	}

	var des []fs.DirEntry
	if !info.IsDir() || *dLsFlag {
		des = append(des, fs.FileInfoToDirEntry(info))
	} else {
		var err error
		des, err = fs.ReadDir(fsys, name)
		if err != nil {
			return err
		}
	}

	if !*lLsFlag {
		for _, de := range des {
			if _, err := fmt.Fprintln(w, de.Name()); err != nil {
				return err
			}
		}
		return nil
	}

	infos := make([]fs.FileInfo, 0, len(des))
	for _, de := range des {
		i, err := de.Info()
		if err != nil {
			return err
		}
		infos = append(infos, i)
	}

	longfmts := LsLongFormats(infos)
	for i := range infos {
		_, err := fmt.Fprintf(w, "%s %s\n", longfmts[i], infos[i].Name())
		if err != nil {
			return err
		}
	}

	return nil
}

func LsLongFormats(infos []fs.FileInfo) []string {
	var modeWidth, sizeWidth int
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

	longfmts := make([]string, 0, len(infos))
	for _, info := range infos {
		t := info.ModTime()
		var timeStr string
		if t.Year() == time.Now().Year() {
			timeStr = t.Format("Jan _2 15:04")
		} else {
			timeStr = t.Format("Jan _2  2006")
		}

		longfmts = append(longfmts, fmt.Sprintf("%-*s %*d %s",
			modeWidth,
			info.Mode().String(),
			sizeWidth,
			info.Size(),
			timeStr))
	}
	return longfmts
}
