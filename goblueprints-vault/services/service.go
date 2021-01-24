package services

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Hash (ctx context.Context,password string)(string,error)
	Validate(ctx context.Context,password,hash string)(bool,error)
}

func NewService()Service{
	return vaultService{}
}

type vaultService struct {
}

func (vaultService) Hash(ctx context.Context, password string) (string, error) {
	fmt.Println(password)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (vaultService) Validate(ctx context.Context, password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, nil
	}
	return true, nil
}
