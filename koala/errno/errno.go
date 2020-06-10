package errno

import "fmt"

type KoalaError struct {
	Code int
	Message string
}

func(k *KoalaError)Error()string{
	return fmt.Sprintf("koaka error,code:%d message:%v",k.Code,k.Message)
}

var (
	NotHaveInstance = &KoalaError{
		Code:    1,
		Message: "not have instance",
	}
	ConnFailed = &KoalaError{
		Code:    2,
		Message: "connect failed",
	}
	InvalidNode = &KoalaError{
		Code:    3,
		Message: "invalid node",
	}
	AllNodeFailed = &KoalaError{
		Code:    4,
		Message: "all node failed",
	}
)

func IsConnectError(err error)(res bool){
	koalaErr,ok := err.(*KoalaError)
	if !ok{
		return
	}
	if koalaErr ==ConnFailed{
		res = true
	}
	return
}