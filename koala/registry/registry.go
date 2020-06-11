package registry

import "context"

type Registry interface {
	Name()string
	Init(ctx context.Context,opts ...Option)error
	Register(ctx context.Context,service *Service)error
	Unregister(ctx context.Context,service *Service)error
	GetService(ctx context.Context,name string)(*Service,error)
}