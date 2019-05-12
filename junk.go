package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	schemalink = `create table if not exists
		xlink (
		dir text,
		nam text,
		dev integer,
		ino integer,
		siz integer,
		mod integer
		)`
	insertlink = `insert into
		xlink (
		dir,
		nam,
		dev,
		ino,
		siz,
		mod
		) values (
		:dir,
		:nam,
		:dev,
		:ino,
		:siz,
		:mod
		)`
)

type xlink struct {
	Dir string `db:"dir"`
	Nam string
	Dev int32
	Ino uint64
	Siz int64
	Mod uint16
}

func mkwf(verbose bool, stmtlink *sqlx.NamedStmt) filepath.WalkFunc {
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

		stmtlink.MustExec(&link)

		return nil
	}
}

func main() {
	var e error

	dbfilename := os.Args[1]
	startpath := os.Args[2]

	db, e := sqlx.Open("sqlite3", "file:"+dbfilename)
	if e != nil {
		panic(e)
	}
	defer db.Close()

	e = db.Ping()
	if e != nil {
		panic(e)
	}

	db.MustExec("drop table if exists xlink")
	db.MustExec(schemalink)

	stmtlink, e := db.PrepareNamed(insertlink)
	if e != nil {
		panic(e)
	}
	defer stmtlink.Close()

	e = filepath.Walk(startpath, mkwf(true, stmtlink))
	if e != nil {
		panic(e)
	}
}
