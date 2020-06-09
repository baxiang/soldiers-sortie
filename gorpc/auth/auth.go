package auth

import (
	"context"
	"net"
)

type AuthInfo interface {
	AuthType() string
}

type TransportAuth interface {
	ClientHandshake(context.Context, string, net.Conn) (net.Conn, AuthInfo, error)
}

// PerRPCAuth defines a common interface for single RPC call authentication
type PerRPCAuth interface {

	// GetMetadata fetch custom metadata from the context
	GetMetadata(ctx context.Context, uri ... string) (map[string]string, error)

}