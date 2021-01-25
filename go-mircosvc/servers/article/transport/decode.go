package transport

import (
	"context"
	"errors"

	articlePb "github.com/baxiang/soldiers-sortie/go-mircosvc/pb"
)

func decodeGRPCGetCategoriesResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	rp, ok := grpcResponse.(*articlePb.GetCategoriesResponse)
	if !ok {
		return nil, errors.New("decodeGRPCGetCategoriesResponse: interface conversion error")
	}

	return rp, nil
}