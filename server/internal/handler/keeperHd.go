package handler

import (
	"context"

	"github.com/boginskiy/GophKeeper/server/internal/rpc"
)

type KeeperHandler struct {
	rpc.UnimplementedKeeperServiceServer
	// Add Service
}

func NewKeeperHandler() *KeeperHandler {
	return &KeeperHandler{}
}

func (k *KeeperHandler) RegistUser(ctx context.Context, req *rpc.RegistUserRequest) (*rpc.RegistUserResponse, error) {
	return &rpc.RegistUserResponse{}, nil
}

func (k *KeeperHandler) AuthUser(ctx context.Context, req *rpc.AuthUserRequest) (*rpc.AuthUserResponse, error) {
	return &rpc.AuthUserResponse{}, nil
}
