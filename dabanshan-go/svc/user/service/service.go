package service

import (
	"context"
	"errors"
	"github.com/baxiang/soldiers-sortie/dabanshan-go/svc/user/model"
)

var (
	// ErrUserNotFound 用户未发现
	ErrUserNotFound = errors.New("not found user")
	// ErrUserAlreadyExisting 用户名已存在
	ErrUserAlreadyExisting = errors.New("username already existing")
)

type Service interface {
	GetUser(ctx context.Context, id string) (model.GetUserResponse, error)
	Register(ctx context.Context, RegisterRequest model.RegisterRequest) (model.RegisterUserResponse, error)
	Login(ctx context.Context, login model.LoginRequest) (model.LoginResponse, error)
	//Upload(ctx context.Context, manifestName string, manifest io.Reader, fileName string, file io.Reader) (string, error)
}