package client

import (
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/intercept"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientGRPC struct {
	AuthService   rpc.AuthServiceClient
	TexterService rpc.TexterServiceClient
	ByterService  rpc.ByterServiceClient
	Conn          *grpc.ClientConn
}

func NewClientGRPC(config config.Config, logger logg.Logger, intrcept intercept.ClientInterceptor) *ClientGRPC {
	// Conn
	conn, err := grpc.NewClient(
		config.GetServerGrpc(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(intrcept.SingleAuth),
		grpc.WithChainStreamInterceptor(intrcept.StreamAuth))

	logger.CheckWithFatal(err, "error created client gRPC")

	return &ClientGRPC{
		AuthService:   rpc.NewAuthServiceClient(conn),
		TexterService: rpc.NewTexterServiceClient(conn),
		ByterService:  rpc.NewByterServiceClient(conn),
		Conn:          conn,
	}
}

func (c *ClientGRPC) Close() error {
	return c.Conn.Close()
}
