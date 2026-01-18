package server

import (
	"net"

	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/internal/intercept"
	"github.com/boginskiy/GophKeeper/server/internal/logg"
	"github.com/boginskiy/GophKeeper/server/internal/rpc"
	"google.golang.org/grpc"
)

type ServerGRPC struct {
	Cfg    config.Config
	Logg   logg.Logger
	listen net.Listener
}

func NewServerGRPC(config config.Config, logger logg.Logger) *ServerGRPC {
	// Port for server.
	lst, err := net.Listen("tcp", config.GetPortServerGRPC())
	logger.CheckWithFatal(err, "server listener initialization error")

	return &ServerGRPC{
		Cfg:    config,
		Logg:   logger,
		listen: lst,
	}
}

func (s *ServerGRPC) Run(handler rpc.KeeperServiceServer, interceptor intercept.Interceptor) {
	// New server.
	server := grpc.NewServer(grpc.UnaryInterceptor(interceptor.WithAuth))
	// Registration handler.
	rpc.RegisterKeeperServiceServer(server, handler)
	// Run.
	s.Logg.CheckWithFatal(server.Serve(s.listen), "gRPC server has not started")
}
