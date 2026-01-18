package intercept

import (
	"context"
	"fmt"

	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/internal/logg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Interceptor interface {
	WithAuth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
}

type Intercept struct {
	Cfg  config.Config
	Logg logg.Logger
}

func NewIntercept(config config.Config, logger logg.Logger) *Intercept {
	return &Intercept{
		Cfg:  config,
		Logg: logger,
	}
}

func (i *Intercept) WithAuth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// И так в контексте передаем токен,
	fmt.Println("RRRRR")
	return nil, status.Error(codes.Unauthenticated, "token is bad")
}
