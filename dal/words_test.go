package dal

import (
	"testing"
	"vocabulary-record/model"
	_ "vocabulary-record/test_helper"
	"vocabulary-record/utils"
)

func TestQueryWords(t *testing.T) {
	result, err := QueryWords(0, 10, WordsConds{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

func TestInsertWords(t *testing.T) {
	err := InsertWords([]model.Words{
		{Word: "grazing", Meaning: nil},
		{Word: "seclusion", Meaning: nil},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateWords(t *testing.T) {
	rowsAffected, err := UpdateWords(model.Words{
		Word:    "grazing",
		Meaning: nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rowsAffected)
}

func TestUpdateWordsFailure(t *testing.T) {
	rowsAffected, err := UpdateWords(model.Words{
		Word:    "hello",
		Meaning: utils.ToStringPtr("放牧"),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rowsAffected)
}
