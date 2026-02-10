package codec

import (
	"github.com/boginskiy/GophKeeper/server/internal/rpc"
)

type AuthGRPCDecoder[T any] interface {
	DecoderRegistRequest(req *rpc.RegistRequest) (T, error)
	DecoderRecovRequest(req *rpc.RecovRequest) (T, error)
	DecoderAuthRequest(req *rpc.AuthRequest) (T, error)
}

type AuthGRPCEncoder[T any] interface {
	EncodeRegistResponse(req T) (T, error)
}

type ByteGRPCCoder[T any] interface {
	ByteGRPCDecoder[T]
	ByteGRPCEncoder[T]
}

type ByteGRPCDecoder[T any] interface {
	DecoderUploadStreamCtx(stream rpc.ByterService_UploadServer) (T, error)
	DecoderUnloadStreamCtx(stream rpc.ByterService_UnloadServer) (T, error)
}

type ByteGRPCEncoder[T any] interface {
	EncodeDeleteBytesResponse(string) (*rpc.DeleteBytesResponse, error)
	EncodeReadAllBytesResponse([]T) (*rpc.ReadAllBytesResponse, error)
	EncodeUploadBytesResponse(T) (*rpc.UploadBytesResponse, error)
	EncodeReadBytesResponse(T) (*rpc.ReadBytesResponse, error)
}
