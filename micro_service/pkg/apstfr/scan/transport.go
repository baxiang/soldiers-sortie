package scan

import (
	"context"
	"userService/pkg/apstfr/apstfrpb"

	"github.com/go-kit/kit/endpoint"
)

func MakePayEndpoint(server apstfrpb.ScanServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return server.Pay(ctx, request.(*apstfrpb.PayRequest))
	}
}
