package infra

type Checker interface {
	CheckTwoString(oneStr, twoStr string) bool
	CheckPassword(userPassword, password string) bool
	CheckTypeFile(pathToFile, typ string) bool
}
