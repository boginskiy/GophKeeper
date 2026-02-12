package model

import (
	"os"
	"time"
)

type Bytes struct {
	Name         string
	Descr        *os.File
	SentSize     string
	ReceivedSize string
	Type         string
	UpdatedAt    time.Time
}

func (b *Bytes) GetFileType() string {
	return b.Type
}

func (b *Bytes) GetFileName() string {
	return b.Name
}
