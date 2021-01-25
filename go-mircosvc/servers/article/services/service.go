package services

import (
	"context"
	"database/sql"

	"github.com/baxiang/soldiers-sortie/go-mircosvc/pkg/validator"

	articlePb "github.com/baxiang/soldiers-sortie/go-mircosvc/pb"
)

type ArticleServicer interface {
	GetCategories(context.Context) (*articlePb.GetCategoriesResponse, error)
}

func NewArticleService(db *sql.DB) ArticleServicer {
	return &ArticleService{
		db,
		validator.NewValidator(),
	}
}

type ArticleService struct {
	mysql     *sql.DB
	validator *validator.Validator
}

func (svc *ArticleService) GetCategories(_ context.Context) (*articlePb.GetCategoriesResponse, error) {
	return &articlePb.GetCategoriesResponse{
		Count: 1,
		Data: []*articlePb.CategoryResponse{
			&articlePb.CategoryResponse{
				Id:   1,
				Name: "11",
			},
		},
	}, nil
}