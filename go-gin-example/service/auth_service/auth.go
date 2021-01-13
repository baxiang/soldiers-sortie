package auth_service

import "github.com/baxiang/soldiers-sortie/go-gin-example/models"

type Auth struct {
	Username string
	Password string
}

func (a *Auth) Check() (bool, error) {
	return models.CheckAuth(a.Username, a.Password)
}