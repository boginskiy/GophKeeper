package intercept

import (
	"context"

	"google.golang.org/grpc"
)

type ServInterceptor interface {
	ServAuth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
	StreamServAuth(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error
}
