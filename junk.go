package main

import (
	"fmt"
	"io"
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
	dir string
	nam string
	dev int32
	ino uint64
	siz int64
	mod uint16
}

func mkwf(stmtlink *sqlx.NamedStmt) filepath.WalkFunc {
	startpathIsSet := false
	startpath := ""
	counter := 0

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
		left := len(startpath) + 1
		right := len(path) - len(nam) - 1
		dir := ""
		if left < right {
			dir = path[left:right]
		}
		link := xlink{dir, nam, st.Dev, st.Ino, st.Size, st.Mode}

		stmtlink.MustExec(&link)

		fmt.Fprintf(os.Stderr, "%x %8x %6o %9d %q %q\n", link.dev, link.ino, link.mod, link.siz, link.dir, link.nam)
		counter++
		if counter > 10 {
			return io.EOF
		}
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

	db.MustExec(schemalink)

	tx := db.MustBegin()
	defer tx.Commit()

	stmtlink, e := tx.PrepareNamed(insertlink)
	if e != nil {
		panic(e)
	}
	defer stmtlink.Close()

	e = filepath.Walk(startpath, mkwf(stmtlink))
	if e == io.EOF {
		return
	}
	if e != nil {
		panic(e)
	}
}
