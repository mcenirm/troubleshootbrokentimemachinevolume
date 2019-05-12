package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

func mkwf(verbose bool, collector *xcollector) filepath.WalkFunc {
	startpathIsSet := false
	startpath := ""

	return func(path string, info os.FileInfo, err error) error {
		if !startpathIsSet {
			startpathIsSet = true
			startpath = path
			return nil
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil
		}

		nam := info.Name()
		st := info.Sys().(*syscall.Stat_t)
		left := len(startpath)
		right := len(path) - len(nam) - 1
		dir := ""
		if left < right {
			dir = path[left:right]
		}
		link := xlink{dir, nam, st.Dev, st.Ino, st.Size, st.Mode}

		if verbose {
			fmt.Fprintf(os.Stderr, "%x %8x %6o %9d %q %q\n", link.Dev, link.Ino, link.Mod, link.Siz, link.Dir, link.Nam)
		}

		collector.collect(&link)

		return nil
	}
}

func main() {
	var e error

	dbfilename := os.Args[1]
	startpath := os.Args[2]

	collector, e := openCollector(dbfilename)
	if e != nil {
		panic(e)
	}
	defer collector.Close()

	e = filepath.Walk(startpath, mkwf(false, collector))
	if e != nil {
		panic(e)
	}
}
