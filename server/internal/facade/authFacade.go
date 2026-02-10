package facade

import (
	"context"

	"github.com/boginskiy/GophKeeper/server/internal/codec"
	"github.com/boginskiy/GophKeeper/server/internal/errs"
	"github.com/boginskiy/GophKeeper/server/internal/handler2"
	"github.com/boginskiy/GophKeeper/server/internal/model"
	"github.com/boginskiy/GophKeeper/server/internal/rpc"
)

type AuthFacade struct {
	rpc.UnimplementedAuthServiceServer
	ErrMapper   errs.ErrMapper
	AuthHandler handler2.AuthHandler[*model.User]
	AuthDecoder codec.AuthGRPCDecoder[*model.User]
}

func NewAuthFacade(
	errMapper errs.ErrMapper,
	authHandler handler2.AuthHandler[*model.User],
	authDecoder codec.AuthGRPCDecoder[*model.User],
) *AuthFacade {

	return &AuthFacade{
		ErrMapper:   errMapper,
		AuthHandler: authHandler,
		AuthDecoder: authDecoder,
	}
}

func (k *AuthFacade) Registration(ctx context.Context, req *rpc.RegistRequest) (*rpc.RegistResponse, error) {
	modUser, err := k.AuthDecoder.DecoderRegistRequest(req)
	if err != nil {
		return nil, k.ErrMapper.Mapping(err)
	}

	status, err := k.AuthHandler.HandleRegistration(ctx, modUser)
	if err != nil {
		return nil, k.ErrMapper.Mapping(err)
	}

	// EncodeRegistResponse
	return &rpc.RegistResponse{Status: status}, nil
}

func (k *AuthFacade) Authentication(ctx context.Context, req *rpc.AuthRequest) (*rpc.AuthResponse, error) {
	modUser, _ := k.AuthDecoder.DecoderAuthRequest(req)

	status, err := k.AuthHandler.HandleAuthentication(ctx, modUser)
	if err != nil {
		return nil, k.ErrMapper.Mapping(err)
	}

	return &rpc.AuthResponse{Status: status}, nil
}

func (k *AuthFacade) Recovery(ctx context.Context, req *rpc.RecovRequest) (*rpc.RecovResponse, error) {
	modUser, err := k.AuthDecoder.DecoderRecovRequest(req)
	if err != nil {
		return nil, k.ErrMapper.Mapping(err)
	}

	status, err := k.AuthHandler.HandleRecovery(ctx, modUser)
	if err != nil {
		return nil, k.ErrMapper.Mapping(err)
	}

	return &rpc.RecovResponse{Status: status}, nil
}
