package utils

import (
	"github.com/jmoiron/sqlx"
	"log"
)

func ScanRows[T any](db *sqlx.DB, sqlStmt string) (
	[]T,
	error,
) {
	rows, err := db.Queryx(sqlStmt)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	results := make([]T, 0)
	for rows.Next() {
		var result T
		err = rows.StructScan(&result)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		results = append(results, result)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return results, nil
}
