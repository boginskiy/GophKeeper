package errs

import "errors"

var (
	ErrUncorrectCredentials = errors.New("uncorrected user credentials")
	ErrEmptyConfigFile      = errors.New("config file is empty")
)
