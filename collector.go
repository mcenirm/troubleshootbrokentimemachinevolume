package main

import (
	"github.com/jmoiron/sqlx"
)

type xcollector struct {
	stmt *sqlx.NamedStmt
}
