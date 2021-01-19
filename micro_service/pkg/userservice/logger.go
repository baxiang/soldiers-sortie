package userservice

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/sirupsen/logrus"
)

func logginMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		logrus.Debugln("request info:", request)
		return next(ctx, request)
	}
}
