package transport

import "time"

type ClientTransportOptions struct {
	Target string
	ServiceName string
	Network string

	Timeout time.Duration
}


type ClientTransportOption func(*ClientTransportOptions)