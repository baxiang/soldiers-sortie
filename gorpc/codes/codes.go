package codes

import "fmt"
const (
	OK = 0
	ServerInternalErrorCode = 100
	ConfigErrorCode = 101
	NetworkNotSupportedErrorCode = 201
	ClientMsgErrorCode = 301
	ClientCertFail = 401
)

const (
	FrameworkError = 1
	BusinuessError = 2
)

var (
	ServerInternalError = NewFrameworkError(ServerInternalErrorCode,"server internal error")
	ConfigError = NewFrameworkError(ConfigErrorCode,"config error")
	NetworkNotSupportedError = NewFrameworkError(NetworkNotSupportedErrorCode,"network type not supported")
	ClientCertFailError = NewFrameworkError(ClientCertFail, "client cert fail")
)

type Error struct {
	Code uint32
	Type int
	Message string
}

const (
	Success = "success"
)
func (e *Error) Error() string {
	if e == nil {
		return Success
	}
	if e.Type == FrameworkError {
		return fmt.Sprintf("type : framework, code : %d, msg : %s",e.Code, e.Message)
	}
	return fmt.Sprintf("type : business, code : %d, msg : %s",e.Code, e.Message)
}

func NewFrameworkError(code uint32, msg string) *Error{
	return &Error{
		Type : FrameworkError,
		Code : code,
		Message : msg,
	}
}
func New(code uint32, msg string) *Error{
	return &Error{
		Type : BusinuessError,
		Code : code,
		Message : msg,
	}
}