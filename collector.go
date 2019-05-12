package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type xcollector struct {
	db   *sqlx.DB
	stmt *sqlx.NamedStmt
}

func openCollector(dbfilename string) (*xcollector, error) {
	var e error

	c := new(xcollector)

	c.db, e = sqlx.Open("sqlite3", "file:"+dbfilename)
	if e != nil {
		return nil, e
	}

	c.db.MustExec("drop table if exists xlink")
	c.db.MustExec(schemalink)

	c.stmt, e = c.db.PrepareNamed(insertlink)
	if e != nil {
		return c, e
	}

	return c, nil
}

func (c *xcollector) Close() error {
	if c.stmt != nil {
		if e := c.stmt.Close(); e != nil {
			return e
		}
	}
	if c.db != nil {
		if e := c.db.Close(); e != nil {
			return e
		}
	}
	return nil
}

func (c *xcollector) collect(link *xlink) {
	c.stmt.MustExec(&link)
}
