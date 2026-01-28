package model

import (
	"os"

	"github.com/boginskiy/GophKeeper/client/internal/errs"
	"github.com/boginskiy/GophKeeper/client/internal/utils"
)

type Bytes struct {
	Name  string
	Descr *os.File
	Size  int64
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
		Name:  name,
		Descr: descr,
		Size:  size,
	}, nil
}
