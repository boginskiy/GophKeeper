package utils

import "github.com/boginskiy/GophKeeper/client/pkg"

type Check struct {
}

func NewCheck() *Check {
	return &Check{}
}

func (d *Check) CheckTwoString(oneStr, twoStr string) bool {
	return oneStr == twoStr
}

func (d *Check) CheckPassword(userPassword, password string) bool {
	return pkg.CompareHashAndPassword(userPassword, password)
}
