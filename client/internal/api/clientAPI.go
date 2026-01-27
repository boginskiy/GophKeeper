package api

// import (
// 	"context"
// 	"time"

// 	"github.com/boginskiy/GophKeeper/client/cmd/client"
// 	"github.com/boginskiy/GophKeeper/client/cmd/config"
// 	"github.com/boginskiy/GophKeeper/client/internal/errs"
// 	"github.com/boginskiy/GophKeeper/client/internal/logg"
// 	"github.com/boginskiy/GophKeeper/client/internal/rpc"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/metadata"
// 	"google.golang.org/grpc/status"
// )

// type ClientAPI struct {
// 	Cfg        config.Config
// 	Logg       logg.Logger
// 	ClientGRPC *client.ClientGRPC
// }

// func NewClientAPI(
// 	cfg config.Config,
// 	logger logg.Logger,
// 	clientgrpc *client.ClientGRPC) *ClientAPI {

// 	return &ClientAPI{
// 		Cfg:        cfg,
// 		Logg:       logger,
// 		ClientGRPC: clientgrpc,
// 	}
// }

// func (c *ClientAPI) RegisterUser(req *rpc.RegistUserRequest, header *metadata.MD) (*rpc.RegistUserResponse, error) {
// 	// Ctx.
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Cfg.GetWaitingTimeResponse()))
// 	defer cancel()

// 	// Retry.
// 	for retry := c.Cfg.GetCountRetryRequest(); retry >= 0; retry-- {
// 		// Call server GRPC.
// 		res, err := c.ClientGRPC.AuthService.RegistUser(ctx, req, grpc.Header(header))

// 		if statusErr, ok := status.FromError(err); ok && statusErr.Code() == codes.DeadlineExceeded {
// 			retry--
// 		} else {
// 			return res, err
// 		}
// 	}

// 	return nil, errs.ErrResponseServer
// }

// func (c *ClientAPI) AutherUser(req *rpc.AuthUserRequest, header *metadata.MD) (*rpc.AuthUserResponse, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Cfg.GetWaitingTimeResponse()))
// 	defer cancel()

// 	for retry := c.Cfg.GetCountRetryRequest(); retry >= 0; retry-- {
// 		res, err := c.ClientGRPC.AuthService.AuthUser(ctx, req, grpc.Header(header))

// 		if statusErr, ok := status.FromError(err); ok && statusErr.Code() == codes.DeadlineExceeded {
// 			retry--
// 		} else {
// 			return res, err
// 		}
// 	}
// 	return nil, errs.ErrResponseServer
// }

// func (c *ClientAPI) Create(ctx context.Context, req *rpc.CreateRequest) (*rpc.CreateResponse, error) {
// 	var header *metadata.MD

// 	res, err := c.ClientGRPC.TexterService.Create(ctx, req, grpc.Header(header))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return res, nil
// }

// //  rpc Create (CreateRequest) returns (CreateResponse);
// //   rpc Read (ReadRequest) returns (ReadResponse);
// //   rpc ReadAll (ReadAllRequest) returns (ReadAllResponse);
// //   rpc Update (CreateRequest) returns (CreateResponse);
// //   rpc Delete (DeleteRequest) returns (DeleteResponse);
