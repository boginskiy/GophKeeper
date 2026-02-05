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
	SingleAuth(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error
	StreamAuth(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error)
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

func (i *ClientIntercept) SingleAuth(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// If there is no registration
	if i.User.GetModelUser() == nil {
		return invoker(ctx, method, req, reply, cc, opts...)
	}

	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.Pairs()
	}

	md.Set("authorization", i.User.GetModelUser().Token)
	newCtx := metadata.NewOutgoingContext(ctx, md)
	return invoker(newCtx, method, req, reply, cc, opts...)
}

func (i *ClientIntercept) StreamAuth(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.Pairs()
	}

	md.Set("authorization", i.User.GetModelUser().Token)
	newCtx := metadata.NewOutgoingContext(ctx, md)
	return streamer(newCtx, desc, cc, method, opts...)
}
