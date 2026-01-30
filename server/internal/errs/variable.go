package errs

import (
	"errors"
)

var (

	// Errors for this package.
	ErrTokenNotValid  = errors.New(`{"error":"token is not valid"}`)
	ErrDataNotValid   = errors.New(`{"error":"data is not valid"}`)
	ErrTokenIsExpired = errors.New(`{"error":"token is expired"}`)

	ErrUserPassword   = errors.New("invalid password")
	ErrTypeConversion = errors.New("object type conversion error")

	ErrDataOwner = errors.New("owner of the data is not defined")

	// RepoUser
	ErrEmailNotUnique = errors.New("email is not unique")
	ErrUserNotFound   = errors.New("user not found")

	// RepoText
	ErrDataNotFound = errors.New("data not found")

	// Errors with wrap.
	ErrCreateUser  = NewErrWrapper("user creation error")
	ErrCreateToken = NewErrWrapper("token creation error")
)
