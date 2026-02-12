package model

import (
	"time"
)

type Text struct {
	Name         string
	Type         string
	Content      string
	Owner        string
	ListActivate []int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewText(name, typ, content, owner string) *Text {
	return &Text{
		Name:         name,
		Type:         typ,
		Content:      content,
		Owner:        owner,
		ListActivate: make([]int64, 10),
	}
}
