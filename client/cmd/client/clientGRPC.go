package client

import (
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type ClientGRPC struct {
	AuthService   rpc.AuthServiceClient
	TexterService rpc.TexterServiceClient
	ByterService  rpc.ByterServiceClient
	Conn          *grpc.ClientConn
}

func NewClientGRPC(config config.Config, logger logg.Logger) *ClientGRPC {
	// Conn
	conn, err := grpc.NewClient(
		config.GetPortServerGRPC(), grpc.WithTransportCredentials(insecure.NewCredentials()))

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

func (c *ClientGRPC) TakeValueFromHeader(header metadata.MD, field string, idx int) string {
	values := header.Get(field)
	if len(values) > 0 {
		return values[idx]
	}
	return ""
}

func (c *ClientGRPC) CreateHeaderWithValue(key, value string) metadata.MD {
	return metadata.Pairs(key, value)
}
