package service

import "context"

type Servicer interface {
	Create(context.Context, any) (any, error)
	Read(context.Context, any) (any, error)
	ReadAll(context.Context, any) (any, error)
	Update(context.Context, any) (any, error)
	Delete(context.Context, any) (any, error)
}
