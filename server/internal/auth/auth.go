package auth

import (
	"context"

	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/internal/logg"
	"github.com/boginskiy/GophKeeper/server/internal/rpc"
	"google.golang.org/grpc/metadata"
)

type Auth struct {
	Cfg        config.Config
	Logg       logg.Logger
	JWTService JWTer
}

func NewAuth(config config.Config, logger logg.Logger, jwtSrv JWTer) *Auth {
	return &Auth{Cfg: config, Logg: logger, JWTService: jwtSrv}
}

func (a *Auth) Identification(ctx context.Context, req any) (*ExtraInfoToken, bool) {
	// Check, if client go to Authentication.
	if a.Authentication(ctx, req) {
		return nil, true
	}
	// Check, if client go to Registration.
	if a.Registration(ctx, req) {
		return nil, true
	}

	token := a.takeDataFromCtx(ctx, "authorization")

	// Try Authentication.
	infoToken, err := a.JWTService.GetInfoAndValidJWT(token)
	if err != nil {
		return nil, false
	}

	return infoToken, true
}

func (a *Auth) Registration(ctx context.Context, req any) bool {
	_, ok := req.(*rpc.RegistUserRequest)
	return ok
}

func (a *Auth) Authentication(ctx context.Context, req any) bool {
	_, ok := req.(*rpc.AuthUserRequest)
	return ok
}

func (a *Auth) takeDataFromCtx(ctx context.Context, data string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		val := md.Get(data)
		if len(val) > 0 {
			return val[0]
		}
	}
	return ""
}
