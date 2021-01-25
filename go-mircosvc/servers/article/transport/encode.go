package transport

import (
	"context"
	"errors"

	articlePb "github.com/baxiang/soldiers-sortie/go-mircosvc/pb"
	"github.com/baxiang/soldiers-sortie/go-mircosvc/pkg/common"
)

func encodeGRPCGetCategoriesResponse(_ context.Context, response interface{}) (interface{}, error) {
	res, ok := response.(common.Response)
	if !ok {
		return nil, errors.New("encodeGRPCGetCategoriesResponse: interface conversion error")
	}

	data := res.Data.(*articlePb.GetCategoriesResponse)

	return data, nil
}
