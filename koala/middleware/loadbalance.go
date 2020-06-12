package middleware

import (
	"context"
	"github.com/baxiang/koala/errno"
	"github.com/baxiang/koala/loadbalance"
	"github.com/baxiang/koala/logs"
	"github.com/baxiang/koala/meta"
)

func NewLoadBalanceMiddleware(balancer loadbalance.LoadBalance) Middleware{
	return func(next MiddlewareFunc) MiddlewareFunc {
		return func(ctx context.Context, req interface{}) (rsp interface{}, err error) {
			rpcMeta := meta.GetRpcMeta(ctx)
			if len(rpcMeta.AllNodes) == 0 {
				err = errno.NotHaveInstance
				logs.Error(ctx, "not have instance")
				return
			}
			ctx = loadbalance.WithBalanceContext(ctx)
			for  {
				rpcMeta.CurrNode, err = balancer.Select(ctx, rpcMeta.AllNodes)
				if err != nil {
					return
				}
				logs.Debug(ctx, "select node:%#v", rpcMeta.CurrNode)
				rpcMeta.HistoryNodes = append(rpcMeta.HistoryNodes, rpcMeta.CurrNode)
				rsp, err = next(ctx, req)
				if err != nil {
					//连接错误的话，进行重试
					if errno.IsConnectError(err) {
						continue
					}
					return
				}
				break
			}
			return
		}

	}
}