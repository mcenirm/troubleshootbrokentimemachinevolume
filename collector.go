package main

import (
	"github.com/jmoiron/sqlx"
)

type xcollector struct {
	stmt *sqlx.NamedStmt
}

func (c *xcollector) collect(link *xlink) {
	c.stmt.MustExec(&link)
}
