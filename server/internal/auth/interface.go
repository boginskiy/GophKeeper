package auth

import "context"

// JWTer for a work of JWT authentication.
type JWTer interface {
	GetInfoAndValidJWT(token string) (*ExtraInfoToken, error)
	CreateJWT(userID int) (string, error)
}

// Auther .
type Auther interface {
	Identification(context.Context, any) (*ExtraInfoToken, bool)
}
