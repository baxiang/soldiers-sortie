package meta

import (
	"context"
	"github.com/baxiang/koala/registry"
	"google.golang.org/grpc"
)

type RpcMeta struct {
	Caller string //调用方名字
	ServiceName string //服务提供方
	Method string
	CallerCluster string

	ServiceCluster string

	TraceID string
	Env string

	CallerIDC string
	ServiceIDC string

	CurrNode *registry.Node
	HistoryNodes []*registry.Node
	AllNodes []*registry.Node

	Conn *grpc.ClientConn
}

type rpcMetaContextKey struct {}

func InitRpcMeta(ctx context.Context,service,method,caller string)context.Context{
	meta :=&RpcMeta{
		Caller:         caller,
		ServiceName:    service,
		Method:         method,
	}
	return context.WithValue(ctx,rpcMetaContextKey{},meta)
}

func GetRpcMeta(ctx context.Context)*RpcMeta{
	meta,ok := ctx.Value(rpcMetaContextKey{}).(*RpcMeta)
	if !ok{
		meta = &RpcMeta{}
	}
	return meta
}