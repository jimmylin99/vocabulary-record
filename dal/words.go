package dal

import (
	"database/sql"
	"fmt"
	"log"
	"vocabulary-record/model"
	"vocabulary-record/utils"
)

type WordsConds struct {
	Word string
}

func QueryWords(offset int, limit int, conds WordsConds) (
	[]model.Words,
	error,
) {
	sqlStmt := "select * from words order by ID desc"
	if conds.Word != "" {
		sqlStmt += fmt.Sprintf(" where word = '%s'", conds.Word)
	}
	sqlStmt += fmt.Sprintf(" limit %d, %d", offset, limit)

	words, err := utils.ScanRows[model.Words](db, sqlStmt)

	return words, err
}

func InsertWords(values []model.Words) error {
	_, err := db.NamedExec(`INSERT INTO words (word, meaning, occurred_cnt, familiar_cnt)
		VALUES (:word, :meaning, :occurred_cnt, :familiar_cnt)`, values)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func UpdateWords(value model.Words) (
	rowsAffected int64,
	err error,
) {
	var result sql.Result
	result, err = db.NamedExec(`UPDATE words 
SET meaning = :meaning,
    occurred_cnt = :occurred_cnt,
    familiar_cnt = :familiar_cnt
WHERE word = :word`, value)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return
}
