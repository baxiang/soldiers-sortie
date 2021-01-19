package camunda

import (
	"userService/pkg/camunda/pb"
	servicepb "userService/pkg/pb"
)

type Error interface {
	GetErr() *pb.CamundaError
	GetCode() int64
}

func CheckError(err Error) bool {
	return err.GetErr() != nil
}

func TransError(err Error) *servicepb.Error {
	return &servicepb.Error{
		Code:        int32(err.GetCode()),
		Message:     err.GetErr().Type,
		Description: err.GetErr().Message,
	}
}
