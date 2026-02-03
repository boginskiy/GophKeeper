package infra

import (
	"github.com/boginskiy/GophKeeper/client/internal/utils"
	"github.com/boginskiy/GophKeeper/client/pkg"
)

var (
	SOUND = map[string]struct{}{"mp3": {}, "wav": {}, "flac": {}, "ogg": {}}
)

type Check struct {
	FileHandler utils.FileHandler
}

func NewCheck(fileHdlr utils.FileHandler) *Check {
	return &Check{FileHandler: fileHdlr}
}

func (d *Check) CheckTwoString(oneStr, twoStr string) bool {
	return oneStr == twoStr
}

func (d *Check) CheckPassword(userPassword, password string) bool {
	return pkg.CompareHashAndPassword(userPassword, password)
}

func (d *Check) checkTypeOfSound(fileName string) bool {
	typeFile := d.FileHandler.GetTypeFile(fileName)
	if _, ok := SOUND[typeFile]; ok {
		return ok
	}
	return false
}

func (d *Check) CheckTypeFile(filePath, typ string) bool {
	fileName := d.FileHandler.TakeFileFromPath(filePath)
	switch typ {
	case "sound":
		return d.checkTypeOfSound(fileName)
	}
	return false
}
