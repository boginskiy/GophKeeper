package main

import (
	"context"
	"fmt"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {

	logger := logg.NewLogg("main.log", "INFO")
	cfg := config.NewConf(logger)

	var header metadata.MD
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "ASDFG")

	clientGRPC := client.NewClientGRPC(cfg, logger)

	// Собираем инфу для регистрации

	_, err := clientGRPC.Service.RegistUser(ctx, &rpc.RegistUserRequest{}, grpc.Header(&header))

	fmt.Println(err)

}
