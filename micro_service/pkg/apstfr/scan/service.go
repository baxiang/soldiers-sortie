package scan

import (
	"context"
	"userService/pkg/apstfr/apstfrpb"
)

type service struct{}

func (s service) Pay(ctx context.Context, in *apstfrpb.PayRequest) (*apstfrpb.PayReply, error) {
	panic("implement me")
}
