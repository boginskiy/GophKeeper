package auth

import (
	"errors"
	"fmt"
	"time"

	conf "github.com/boginskiy/Clicki/cmd/config"
	"github.com/golang-jwt/jwt/v4"
)

type ExtraInfoToken struct {
	Email       string
	PhoneNumber string
}

// Claims - own statement.
type Claims struct {
	jwt.RegisteredClaims
	InfoToken *ExtraInfoToken
}

// JWTService - is JWT authentication.
type JWTService struct {
	Cfg conf.Config
}

func NewJWTService(config conf.Config) *JWTService {
	return &JWTService{Cfg: config}
}

// CreateToken .
func (j *JWTService) CreateJWT(infoToken *ExtraInfoToken) (string, error) {
	// New jwtToken.
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(time.Duration(j.Cfg.GetTokenLiveTime())))},
		InfoToken: infoToken,
	})

	// Line of token.
	token, err := jwtToken.SignedString([]byte(j.Cfg.GetSecretKey()))
	if err != nil {
		return "", err
	}
	return token, nil
}

// GetInfoAndValidJWT .
func (j *JWTService) GetInfoAndValidJWT(checkingToken string) (*ExtraInfoToken, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(checkingToken, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(j.Cfg.GetSecretKey()), nil
	})

	if err != nil {
		// Analize of expired token.
		if j.checkTokenIsExpired(err) {
			return claims.InfoToken, ErrTokenIsExpired
		}
		return claims.InfoToken, err
	}

	// Analize of invalid token.
	if !token.Valid {
		return claims.InfoToken, ErrTokenNotValid
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
