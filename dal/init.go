package dal

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sqlx.DB
)

func init() {
	var err error
	db, err = sqlx.Connect("sqlite3", "model/vocabulary")
	if err != nil {
		panic(err)
	}
}
