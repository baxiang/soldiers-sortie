package util

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"golang.org/x/time/rate"
)

var ErrLimited = errors.New("rate limit exceeded")

func NewLimiterMiddleware(bucket *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if bucket.Allow() {
				return next(ctx, request)
			}
			return nil, ErrLimited
		}
	}
}
