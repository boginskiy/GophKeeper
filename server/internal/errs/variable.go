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
	ErrDataCtx        = errors.New("context data error")

	// RepoUser
	ErrEmailNotUnique = errors.New("email is not unique")
	ErrUserNotFound   = errors.New("user not found")
	ErrFileNotFound   = errors.New("file not found")

	// RepoText
	ErrDataNotFound = errors.New("data not found")

	// Errors with wrap.
	ErrCreateUser           = NewErrWrapper("user creation error")
	ErrCreateToken          = NewErrWrapper("token creation error")
	ErrCreatePathToStore    = NewErrWrapper("store path creation error")
	ErrCreateFoldersToStore = NewErrWrapper("store folder creation error")
	ErrCreateFile           = NewErrWrapper("file creation error")
	ErrRunStream            = NewErrWrapper("run stream error")

	ErrReadFileToBuff = NewErrWrapper("error reading the file into the buffer")
	ErrSendChankFile  = NewErrWrapper("error sending the part of file to server")

	// DB
	ErrPingDataBase = errors.New("bad database ping")
)
