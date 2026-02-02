package model

import (
	"os"
	"strconv"
	"time"

	"github.com/boginskiy/GophKeeper/client/internal/errs"
	"github.com/boginskiy/GophKeeper/client/internal/utils"
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

func NewBytesFromFile(fileHdlr utils.FileHandler, pathToFile string) (*Bytes, error) {
	pathToFile, err := fileHdlr.TransPathToAbs(pathToFile)

	if err != nil || !fileHdlr.CheckOfFile(pathToFile) {
		return nil, errs.ErrReadFile.Wrap(err)
	}

	name := fileHdlr.TakeFileFromPath(pathToFile)
	descr, err := fileHdlr.TakeDescrFromFile(pathToFile)
	if err != nil {
		return nil, err
	}

	size, err := fileHdlr.TakeSizeFile(descr)
	if err != nil {
		return nil, err
	}

	return &Bytes{
		Name:     name,
		Descr:    descr,
		SentSize: strconv.FormatInt(size, 10),
	}, nil
}
