package gorpc

import "context"

type Handler func(interface{},context.Context,func(interface{})error)(interface{},error)
type Service interface {
	Register(string,Handler)
}