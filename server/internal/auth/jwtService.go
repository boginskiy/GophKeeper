package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/internal/errs"
	"github.com/golang-jwt/jwt/v4"
)

type ExtraInfoToken struct {
	Email       string
	PhoneNumber string
}

func NewExtraInfoToken(email, phone string) *ExtraInfoToken {
	return &ExtraInfoToken{Email: email, PhoneNumber: phone}
}

// Claims - own statement.
type Claims struct {
	jwt.RegisteredClaims
	InfoToken *ExtraInfoToken
}

// JWTService - is JWT authentication.
type JWTService struct {
	Cfg config.Config
}

func NewJWTService(conf config.Config) *JWTService {
	return &JWTService{Cfg: conf}
}

// CreateToken .
func (j *JWTService) CreateJWT(infoToken *ExtraInfoToken) (string, error) {
	// New jwtToken.
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(time.Duration(j.Cfg.GetTokenLifetime())))},
		InfoToken: infoToken,
	})

	// Line of token.
	token, err := jwtToken.SignedString([]byte(j.Cfg.GetTokenSecretKey()))
	if err != nil {
		return "", err
	}
	return token, nil
}

// .
func (j *JWTService) GetInfoAndValidJWT(checkingToken string) (*ExtraInfoToken, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(checkingToken, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(j.Cfg.GetTokenSecretKey()), nil
	})

	if err != nil {
		// Analize of expired token.
		if j.checkTokenIsExpired(err) {
			return claims.InfoToken, errs.ErrTokenIsExpired
		}
		return claims.InfoToken, err
	}

	// Analize of invalid token.
	if !token.Valid {
		return claims.InfoToken, errs.ErrTokenNotValid
	}
	return claims.InfoToken, nil
}

// checkTokenIsExpired Analize of expired token.
func (j *JWTService) checkTokenIsExpired(err error) bool {
	var validErr *jwt.ValidationError

	if errors.As(err, &validErr) {
		// Bit И. Check flag expired token.
		return validErr.Errors&jwt.ValidationErrorExpired != 0
	}
	return false
}
