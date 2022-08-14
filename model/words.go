package model

type Words struct {
	ID          int     `db:"ID"`
	Word        string  `db:"word"`
	Meaning     *string `db:"meaning"`
	OccurredCnt int     `db:"occurred_cnt"`
	FamiliarCnt int     `db:"familiar_cnt"`
}
