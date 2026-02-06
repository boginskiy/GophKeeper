package model

import (
	"time"

	"github.com/boginskiy/GophKeeper/server/internal/rpc"
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

func NewText(req *rpc.CreateRequest) *Text {
	return &Text{
		Name:         req.Name,
		Type:         req.Type,
		Content:      req.Text,
		Owner:        req.Owner,
		ListActivate: req.ListActivate,
	}
}
