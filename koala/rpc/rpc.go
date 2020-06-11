package rpc

import (
	"context"
	"github.com/baxiang/koala/meta"
)

func InitRpcMeta(ctx context.Context,service,method,caller string)context.Context{
	return meta.InitRpcMeta(ctx,service,method,caller)
}