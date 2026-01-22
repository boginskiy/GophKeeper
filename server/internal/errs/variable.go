package errs

import (
	"errors"
)

var (

	// Errors for this package.
	ErrTokenNotValid  = errors.New(`{"error":"token is not valid"}`)
	ErrDataNotValid   = errors.New(`{"error":"data is not valid"}`)
	ErrTokenIsExpired = errors.New(`{"error":"token is expired"}`)

	ErrUserPassword = errors.New("invalid password")

	// RepoUser
	ErrEmailNotUnique = errors.New("email is not unique")
	ErrUserNotFound   = errors.New("user not found")

	// Errors with wrap.
	ErrCreateUser  = NewErrWrapper("user creation error")
	ErrCreateToken = NewErrWrapper("token creation error")
)
