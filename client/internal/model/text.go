package model

import (
	"time"
)

type Text struct {
	Name         string
	Type         string
	Tx           string
	Owner        string
	ListActivate []int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewText(name, tp, tx, owner string) *Text {
	return &Text{
		Name:         name,
		Type:         tp,
		Tx:           tx,
		Owner:        owner,
		ListActivate: make([]int64, 10),
	}
}
