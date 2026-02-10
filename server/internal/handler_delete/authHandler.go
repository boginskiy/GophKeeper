package handlerdelete

// import (
// 	"context"

// 	"github.com/boginskiy/GophKeeper/server/internal/auth"
// 	"github.com/boginskiy/GophKeeper/server/internal/errs"
// 	"github.com/boginskiy/GophKeeper/server/internal/infra"
// 	"github.com/boginskiy/GophKeeper/server/internal/rpc"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// )

// type AuthHandler struct {
// 	rpc.UnimplementedAuthServiceServer
// 	Auth auth.Auther
// }

// func NewAuthHandler(auther auth.Auther) *AuthHandler {
// 	return &AuthHandler{Auth: auther}
// }

// func (k *AuthHandler) Registration(ctx context.Context, req *rpc.RegistRequest) (*rpc.RegistResponse, error) {
// 	token, err := k.Auth.Registration(ctx, req)

// 	// Ошибка создания пользователя.
// 	if err == errs.ErrCreateUser {
// 		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
// 	}

// 	// Ошибка уникального email.
// 	if err == errs.ErrEmailNotUnique {
// 		return nil, status.Errorf(codes.AlreadyExists, "%s", err)
// 	}

// 	// Ошибки сервера.
// 	if err == errs.ErrCreateToken || err != nil {
// 		return nil, status.Errorf(codes.Internal, "%s", err)
// 	}

// 	// Кладем токен в заголовок
// 	err = infra.PutDataToCtx(ctx, "authorization", token)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "%s", err)
// 	}

// 	return &rpc.RegistResponse{Status: "ok"}, nil
// }

// func (k *AuthHandler) Authentication(ctx context.Context, req *rpc.AuthRequest) (*rpc.AuthResponse, error) {
// 	token, err := k.Auth.Authentication(ctx, req)

// 	// Ошибка. Пользователь с таким email не найден.
// 	if err == errs.ErrUserNotFound {
// 		return nil, status.Errorf(codes.NotFound, "%s", err)
// 	}

// 	// Ошибка. Неверный пароль.
// 	if err == errs.ErrUserPassword {
// 		return nil, status.Errorf(codes.Unauthenticated, "%s", err)
// 	}

// 	// Ошибки сервера.
// 	if err == errs.ErrCreateToken || err != nil {
// 		return nil, status.Errorf(codes.Internal, "%s", err)
// 	}

// 	// Кладем токен в заголовок
// 	err = infra.PutDataToCtx(ctx, "authorization", token)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "%s", err)
// 	}

// 	return &rpc.AuthResponse{Status: "ok"}, nil
// }

// func (k *AuthHandler) Recovery(ctx context.Context, req *rpc.RecovRequest) (*rpc.RecovResponse, error) {
// 	token, err := k.Auth.Recovery(ctx, req)

// 	// Ошибка. Пользователь с таким email не найден.
// 	if err == errs.ErrUpdateUser {
// 		return nil, status.Errorf(codes.NotFound, "%s", err)
// 	}

// 	// Ошибки сервера.
// 	if err == errs.ErrCreateToken || err != nil {
// 		return nil, status.Errorf(codes.Internal, "%s", err)
// 	}

// 	// Кладем токен в заголовок
// 	err = infra.PutDataToCtx(ctx, "authorization", token)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "%s", err)
// 	}

// 	return &rpc.RecovResponse{Status: "ok"}, nil
// }
