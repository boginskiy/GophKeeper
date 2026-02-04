package errs

import "errors"

var (
	ErrUncorrectCredentials = errors.New("uncorrected user credentials")
	ErrEmptyConfigFile      = errors.New("config file is empty")
	ErrResponseServer       = errors.New("server response time limit exceeded")
	ErrTypeConversion       = errors.New("object type conversion error")
	ErrFileNotFound         = errors.New("file not found")

	ErrReadFile       = NewErrWrapper("file not found")
	ErrStartStream    = NewErrWrapper("start stream error")
	ErrReadFileToBuff = NewErrWrapper("error reading the file into the buffer")
	ErrSendChankFile  = NewErrWrapper("error sending the part of file to server")
)
