package main

import (
	"context"
	"fmt"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/auth"
	"github.com/boginskiy/GophKeeper/client/internal/clients"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/model"
)

func main() {

	logger := logg.NewLogg("main.log", "INFO")
	cfg := config.NewConf(logger)

	// var header metadata.MD
	// ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "ASDFG")
	// clientGRPC := client.NewClientGRPC(cfg, logger)

	// Собираем инфу для регистрации
	// _, err := clientGRPC.Service.RegistUser(ctx, &rpc.RegistUserRequest{}, grpc.Header(&header))
	ctx := context.Background()
	config := config.NewConf(logger)
	userChan := make(chan *model.User, 1)

	clientGRPC := client.NewClientGRPC(cfg, logger)
	clientAPI := clients.NewClientAPI(cfg, logger, clientGRPC)

	client := auth.NewAuthRemote(ctx, config, logger, userChan, clientAPI)

	user := &model.User{
		UserName:    "Dima",
		Email:       "email",
		PhoneNumber: "8980",
		Password:    "1234",
	}

	client.RemoteRegistration(user)

	fmt.Println(user.StatusError, user.Token)

}
