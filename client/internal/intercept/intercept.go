package intercept

import (
	"context"

	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ClientInterceptor interface {
	ClientAuth(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error
}

type ClientIntercept struct {
	Cfg  config.Config
	Logg logg.Logger
	User user.User
}

func NewClientIntercept(config config.Config, logger logg.Logger, user user.User) *ClientIntercept {
	return &ClientIntercept{
		Cfg:  config,
		Logg: logger,
		User: user,
	}
}

func (i *ClientIntercept) ClientAuth(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.Pairs()
	}

	md.Set("authorization", i.User.GetModelUser().Token)
	newCtx := metadata.NewOutgoingContext(ctx, md)
	return invoker(newCtx, method, req, reply, cc, opts...)
}
