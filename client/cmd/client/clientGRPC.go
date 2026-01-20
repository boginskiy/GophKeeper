package client

import (
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientGRPC struct {
	Service rpc.KeeperServiceClient
	Conn    *grpc.ClientConn
}

func NewClientGRPC(config config.Config, logger logg.Logger) *ClientGRPC {
	// Conn
	conn, err := grpc.NewClient(
		config.GetPortServerGRPC(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	logger.CheckWithFatal(err, "error created client gRPC")

	return &ClientGRPC{
		Service: rpc.NewKeeperServiceClient(conn),
		Conn:    conn,
	}
}

func (c *ClientGRPC) Close() error {
	return c.Conn.Close()
}
