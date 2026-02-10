package codec

import (
	"github.com/boginskiy/GophKeeper/server/internal/model"
	"github.com/boginskiy/GophKeeper/server/internal/rpc"
	"github.com/boginskiy/GophKeeper/server/pkg"
)

type AuthDecoderEncoder struct{}

func NewAuthDecoderEncoder() *AuthDecoderEncoder {
	return &AuthDecoderEncoder{}
}

// Registration Codec

// DecoderRegistRequest
func (a *AuthDecoderEncoder) DecoderRegistRequest(req *rpc.RegistRequest) (*model.User, error) {
	hash, err := pkg.GenerateHash(req.GetPassword())
	if err != nil {
		return nil, err
	}
	return &model.User{
		PhoneNumber: req.GetPhonenumber(),
		UserName:    req.GetUsername(),
		Email:       req.GetEmail(),
		Password:    hash,
	}, nil
}

// EncodeRegistResponse
func (a *AuthDecoderEncoder) EncodeRegistResponse(req *model.User) (*rpc.RegistResponse, error) {
	return nil, nil
}

// Authentication Codec

// DecoderAuthRequest
func (a *AuthDecoderEncoder) DecoderAuthRequest(req *rpc.AuthRequest) (*model.User, error) {
	return &model.User{
		Password: req.GetPassword(),
		Email:    req.GetEmail(),
	}, nil
}

// EncodeAuthResponse
func (a *AuthDecoderEncoder) EncodeAuthResponse(req *model.User) (*rpc.AuthResponse, error) {
	return nil, nil
}

// Recovery Codec

// DecoderRegistRequest
func (a *AuthDecoderEncoder) DecoderRecovRequest(req *rpc.RecovRequest) (*model.User, error) {
	hash, err := pkg.GenerateHash(req.GetPassword())
	if err != nil {
		return nil, err
	}
	return &model.User{
		Email:    req.GetEmail(),
		Password: hash,
	}, nil
}
