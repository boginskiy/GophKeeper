package infra

// ManageContext

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Email struct{}
type Phone struct{}
type ID struct{}

var EmailCtx Email
var PhoneCtx Phone
var IDCtx ID

var (
	// Error.
	ErrOwnerData = errors.New("owner of the data is not defined")
	ErrMetaData  = errors.New("metadata not found")
)

// keycontext - is type of key for values for context request.
type keycontext struct{}

func TakeServerValStrFromCtx(ctx context.Context, key Email) (string, error) {
	value, ok := ctx.Value(key).(string)
	if !ok {
		return "", ErrOwnerData
	}
	return value, nil
}

func TakeServerValInt64FromCtx(ctx context.Context, key ID) (int64, error) {
	value, ok := ctx.Value(key).(int64)
	if !ok {
		return 0, ErrOwnerData
	}
	return value, nil
}

func TakeClientValueFromCtx(ctx context.Context, key string, i int) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ErrMetaData
	}
	return md.Get(key)[i], nil
}

func TakeDataFromCtx(ctx context.Context, data string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		val := md.Get(data)
		if len(val) > 0 {
			return val[0]
		}
	}
	return ""
}

func PutDataToCtx(ctx context.Context, key, val string) error {
	return grpc.SetHeader(ctx, metadata.Pairs(key, val))
}
