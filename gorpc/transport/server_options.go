package transport

import (
	"context"
	"time"
)

type ServerTransportOptions struct {
	Address string
	Network string
	Protocol string
	Timeout time.Duration
	Handler Handler
	SerializationType string
	KeepAlivePeriod time.Duration
}

type Handler interface {
	Handle(context.Context,[]byte)([]byte,error)
}

type ServerTransportOption func(*ServerTransportOptions)

// WithServerAddress returns a ServerTransportOption which sets the value for address
func WithServerAddress(address string) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.Address = address
	}
}

// WithServerNetwork returns a ServerTransportOption which sets the value for network
func WithServerNetwork(network string) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.Network = network
	}
}

// WithProtocol returns a ServerTransportOption which sets the value for protocol
func WithProtocol(protocol string) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.Protocol = protocol
	}
}

// WithServerTimeout returns a ServerTransportOption which sets the value for timeout
func WithServerTimeout(timeout time.Duration) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.Timeout = timeout
	}
}

// WithHandler returns a ServerTransportOption which sets the value for handler
func WithHandler(handler Handler) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.Handler = handler
	}
}

// WithSerialization returns a ServerTransportOption which sets the value for serialization
func WithSerializationType(serializationType string) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.SerializationType = serializationType
	}
}

// WithKeepAlivePeriod returns a ServerTransportOption which sets the value for keepAlivePeriod
func WithKeepAlivePeriod(keepAlivePeriod time.Duration) ServerTransportOption {
	return func(o *ServerTransportOptions) {
		o.KeepAlivePeriod = keepAlivePeriod
	}
}