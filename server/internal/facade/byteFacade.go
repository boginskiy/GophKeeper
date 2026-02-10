package facade

import (
	"context"

	"github.com/boginskiy/GophKeeper/server/internal/codec"
	"github.com/boginskiy/GophKeeper/server/internal/errs"
	"github.com/boginskiy/GophKeeper/server/internal/handler2"
	"github.com/boginskiy/GophKeeper/server/internal/model"
	"github.com/boginskiy/GophKeeper/server/internal/rpc"
	"github.com/boginskiy/GophKeeper/server/internal/service"
)

type ByteFacade struct {
	rpc.UnimplementedByterServiceServer
	ErrMapper     errs.ErrMapper
	ByteCoder     codec.ByteGRPCCoder[*model.Bytes]
	ByteHandler   handler2.ByteHandler[*model.Bytes]
	UploadService service.LoadServicer[rpc.ByterService_UploadServer, *model.Bytes]
	UnloadService service.LoadServicer[rpc.ByterService_UnloadServer, *model.Bytes]
}

func NewByteFacade(
	errMapper errs.ErrMapper,
	byteCoder codec.ByteGRPCCoder[*model.Bytes],
	byteHandler handler2.ByteHandler[*model.Bytes],
	uploadService service.LoadServicer[rpc.ByterService_UploadServer, *model.Bytes],
	unloadService service.LoadServicer[rpc.ByterService_UnloadServer, *model.Bytes],
) *ByteFacade {

	return &ByteFacade{
		ErrMapper:     errMapper,
		ByteCoder:     byteCoder,
		ByteHandler:   byteHandler,
		UploadService: uploadService,
		UnloadService: unloadService,
	}
}

func (k *ByteFacade) Upload(stream rpc.ByterService_UploadServer) error {
	modBytes, err := k.ByteCoder.DecoderUploadStreamCtx(stream)
	if err != nil {
		return k.ErrMapper.Mapping(err)
	}

	modBytes, err = k.UploadService.Load(stream, modBytes)
	if err != nil {
		return k.ErrMapper.Mapping(err)
	}

	response, _ := k.ByteCoder.EncodeUploadBytesResponse(modBytes)

	return k.ErrMapper.Mapping(stream.SendAndClose(response))
}

func (k *ByteFacade) Unload(req *rpc.UnloadBytesRequest, stream rpc.ByterService_UnloadServer) error {
	modBytes, err := k.ByteCoder.DecoderUnloadStreamCtx(stream)
	if err != nil {
		return k.ErrMapper.Mapping(err)
	}

	_, err = k.UnloadService.Load(stream, modBytes)
	return k.ErrMapper.Mapping(err)
}

func (k *ByteFacade) Read(ctx context.Context, req *rpc.ReadBytesRequest) (*rpc.ReadBytesResponse, error) {
	modBytes, err := k.ByteHandler.HandleRead(ctx, nil)
	if err != nil {
		return nil, k.ErrMapper.Mapping(err)
	}

	response, _ := k.ByteCoder.EncodeReadBytesResponse(modBytes)
	return response, nil
}

func (k *ByteFacade) ReadAll(ctx context.Context, req *rpc.ReadAllBytesRequest) (*rpc.ReadAllBytesResponse, error) {
	modBytes, err := k.ByteHandler.HandleReadAll(ctx, nil)
	if err != nil {
		return nil, k.ErrMapper.Mapping(err)
	}

	response, _ := k.ByteCoder.EncodeReadAllBytesResponse(modBytes)
	return response, nil
}

func (k *ByteFacade) Delete(ctx context.Context, req *rpc.DeleteBytesRequest) (*rpc.DeleteBytesResponse, error) {
	status, err := k.ByteHandler.HandleDelete(ctx, nil)
	if err != nil {
		return nil, k.ErrMapper.Mapping(err)
	}

	response, err := k.ByteCoder.EncodeDeleteBytesResponse(status)
	return response, nil
}
