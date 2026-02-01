package intercept

import (
	"context"

	"google.golang.org/grpc"
)

// WrapServerStream is Wrap for ServerStream
type WrapServerStream struct {
	grpc.ServerStream
	Ctx context.Context
}

// Переопределяем метод Context, чтобы вернуть наш новый контекст
func (w *WrapServerStream) Context() context.Context {
	return w.Ctx
}
