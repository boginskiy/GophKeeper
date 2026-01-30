package model

import "os"

type Bytes struct {
	Name  string
	Descr *os.File
	Size  int64
}
