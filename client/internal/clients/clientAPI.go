package clients

import (
	"context"
	"time"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/errs"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ClientAPI struct {
	Cfg        config.Config
	Logg       logg.Logger
	ClientGRPC *client.ClientGRPC
}

func NewClientAPI(
	cfg config.Config,
	logger logg.Logger,
	clientgrpc *client.ClientGRPC) *ClientAPI {

	return &ClientAPI{
		Cfg:        cfg,
		Logg:       logger,
		ClientGRPC: clientgrpc,
	}
}

func (c *ClientAPI) RegisterUser(req *rpc.RegistUserRequest, header *metadata.MD) (*rpc.RegistUserResponse, error) {
	// Ctx.
	ctxGRPC, cancel := context.WithTimeout(context.Background(), time.Duration(c.Cfg.GetWaitingTimeResponse()))
	defer cancel()

	// Retry.
	for retry := c.Cfg.GetCountRetryRequest(); retry >= 0; retry-- {
		// Call server GRPC.
		res, err := c.ClientGRPC.Service.RegistUser(ctxGRPC, req, grpc.Header(header))

		if statusErr, ok := status.FromError(err); ok && statusErr.Code() == codes.DeadlineExceeded {
			retry--
		} else {
			return res, err
		}
	}
	return nil, errs.ErrResponseServer
}

func (c *ClientAPI) TakeValueFromHeader(header metadata.MD, field string, idx int) string {
	values := header.Get(field)
	if len(values) > 0 {
		return values[idx]
	}
	return ""
}
