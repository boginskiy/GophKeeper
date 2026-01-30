package server

import (
	"fmt"
	"net"
	"os"

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
	S      *grpc.Server
}

func NewServerGRPC(config config.Config, logger logg.Logger, interceptor intercept.Interceptor) *ServerGRPC {
	// Port for server.
	lst, err := net.Listen("tcp", config.GetPortServerGRPC())
	logger.CheckWithFatal(err, "server listener initialization error")

	server := grpc.NewServer(grpc.UnaryInterceptor(interceptor.WithAuth))

	return &ServerGRPC{
		Cfg:    config,
		Logg:   logger,
		listen: lst,
		S:      server,
	}
}

func (s *ServerGRPC) Registration(
	authSrv rpc.AuthServiceServer,
	texterSrv rpc.TexterServiceServer,
	byterSrv rpc.ByterServiceServer) {

	// Registration services.
	rpc.RegisterTexterServiceServer(s.S, texterSrv)
	rpc.RegisterByterServiceServer(s.S, byterSrv)
	rpc.RegisterAuthServiceServer(s.S, authSrv)
}

func (s *ServerGRPC) Run() {
	fmt.Fprintf(os.Stdout, "Protocol:      %s\n", "gRPC")
	s.Logg.CheckWithFatal(s.S.Serve(s.listen), "gRPC server has not started")
}
