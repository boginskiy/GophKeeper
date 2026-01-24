package model

import "time"

type Meta struct {
	Owner        string
	ListActivate []int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Text struct {
	Name string
	Tx   string
	Mt   *Meta
}

func NewText(name, tx, owner string) *Text {
	return &Text{
		Name: name,
		Tx:   tx,
		Mt: &Meta{
			Owner:     owner,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}
