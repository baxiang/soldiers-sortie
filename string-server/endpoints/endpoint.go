package endpoints

import (
	"context"
	"errors"
	"fmt"

	"github.com/baxiang/soldiers-sortie/string-server/services"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/log"
)

type Endpoints struct {
	GetIsPalindrome endpoint.Endpoint
	GetReverse      endpoint.Endpoint
}

func MakeEndpoints(svc services.StrServiceImp, logger log.Logger) Endpoints {
	return Endpoints{
		GetIsPalindrome: makeGetIsPalindromeEndpoint(svc, logger),
		GetReverse:      makeGetReverseEndpoint(svc, logger),
	}
}

func makeGetIsPalindromeEndpoint(svc services.StrServiceImp, logger log.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println(request)
		req, ok := request.(*IsPalRequest)
		if !ok {
			level.Error(logger).Log("message", "invalid request")
			return nil, errors.New("invalid request")
		}
		msg := svc.IsPal(ctx, req.Word)
		return &IsPalResponse{
			Message: msg,
		}, nil
	}
}


func makeGetReverseEndpoint(svc services.StrServiceImp, logger log.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*ReverseRequest)
		if !ok {
			level.Error(logger).Log("message", "invalid request")
			return nil, errors.New("invalid request")
		}
		reverseString := svc.Reverse(ctx, req.Word)
		return &ReverseResponse{
			Word: reverseString,
		}, nil
	}
}
