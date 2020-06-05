package interceptor

import (
	"context"
)

type Handler func(ctx context.Context,req interface{})(interface{},error)
type ServerInterceptor func(ctx context.Context,req interface{},handler Handler)