package institutionservice

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrDownloadFileNameEmpty = status.Error(codes.NotFound, "empty file name")
	ErrReplyTypeInvalid      = status.Error(codes.Internal, "reply type invalid")
	ErrRequestTypeInvalid    = status.Error(codes.Internal, "request type invalid")
)

const (
	InvalidParam  = "InvalidParamError"
	AlreadyExists = "AlreadyExistsError"
	NotFound      = "NotFoundError"
)
