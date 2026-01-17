package auth

import "errors"

var (
	ErrUncorrectCredentials = errors.New("uncorrected user credentials")
)
