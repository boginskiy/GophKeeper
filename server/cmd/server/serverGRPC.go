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

func NewServerGRPC(config config.Config, logger logg.Logger, intrcep intercept.ServInterceptor) *ServerGRPC {
	// Port for server.
	lst, err := net.Listen("tcp", config.GetServerGrpc())
	logger.CheckWithFatal(err, "server listener initialization error")

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(intrcep.ServAuth),        // Обрабатываем single запросы
		grpc.StreamInterceptor(intrcep.StreamServAuth), // Обрабатываем stream запросы
	}

	server := grpc.NewServer(opts...)

	return &ServerGRPC{
		Cfg:    config,
		Logg:   logger,
		listen: lst,
		S:      server,
	}
}

func (s *ServerGRPC) Registration(
	authService rpc.AuthServiceServer,
	textService rpc.TexterServiceServer,
	bytesService rpc.ByterServiceServer) {

	// Registration services.
	rpc.RegisterTexterServiceServer(s.S, textService)
	rpc.RegisterByterServiceServer(s.S, bytesService)
	rpc.RegisterAuthServiceServer(s.S, authService)
}

func (s *ServerGRPC) Run() {
	fmt.Fprintf(os.Stdout, "Protocol:      %s\n", "gRPC")
	s.Logg.CheckWithFatal(s.S.Serve(s.listen), "gRPC server has not started")
}
