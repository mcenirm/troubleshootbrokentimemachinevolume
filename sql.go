package main

import (
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
