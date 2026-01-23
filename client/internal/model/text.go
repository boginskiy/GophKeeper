package model

import "time"

type Meta struct {
	Owner        string
	ListActivate []int
	CreatedAt    time.Time
	LastActivity time.Time
}

type Text struct {
	Type string
	Name string
	Tx   string
	Mt   Meta
}
