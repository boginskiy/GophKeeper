package infra

import (
	"github.com/boginskiy/GophKeeper/client/internal/utils"
	"github.com/boginskiy/GophKeeper/client/pkg"
)

var (
	IMAGE = map[string]struct{}{"jpg": {}, "jpeg": {}, "png": {}, "gif": {}, "bmp": {}, "webp": {}, "svg": {}, "tif": {}, "tiff": {}, "raw": {}, "ico": {}, "psd": {}}
	VIDEO = map[string]struct{}{"mp4": {}, "avi": {}, "mov": {}, "wmv": {}, "mkv": {}, "webm": {}, "flv": {}}
	SOUND = map[string]struct{}{"mp3": {}, "wav": {}, "flac": {}, "ogg": {}}
	BYTES = map[string]struct{}{}
)

type Check struct {
	FileHandler utils.FileHandler
}

func NewCheck(fileHandler utils.FileHandler) *Check {
	return &Check{FileHandler: fileHandler}
}

func (d *Check) CheckTwoString(oneStr, twoStr string) bool {
	return oneStr == twoStr
}

func (d *Check) CheckPassword(userPassword, password string) bool {
	return pkg.CompareHashAndPassword(userPassword, password)
}

func (d *Check) checkTypeOfMedia(fileName string, listMedia map[string]struct{}) bool {
	typeFile := d.FileHandler.GetTypeFile(fileName)
	if _, ok := listMedia[typeFile]; ok {
		return ok
	}
	return false
}

func (d *Check) CheckTypeFile(filePath, typ string) bool {
	fileName := d.FileHandler.GetFileFromPath(filePath)
	switch typ {
	case "sound":
		return d.checkTypeOfMedia(fileName, SOUND)
	case "video":
		return d.checkTypeOfMedia(fileName, VIDEO)
	case "image":
		return d.checkTypeOfMedia(fileName, IMAGE)
	case "bytes":
		return true
		// return d.checkTypeOfMedia(fileName, BYTES)
	}
	return false
}
