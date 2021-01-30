package services

import (
	"context"
	"strings"
)

type StrServiceImp interface {
	IsPal(context.Context,string) bool     //是否回文串
	Reverse(context.Context,string) string // 翻转字符串
}

type StrService struct {

}

func (svc *StrService) IsPal(ctx context.Context,s string) bool {
	reverse := svc.Reverse(ctx,s)
	return strings.ToLower(s) != reverse
}

func (svc *StrService) Reverse(ctx context.Context,s string) string {
	rns := []rune(s) // convert to rune
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {
		rns[i], rns[j] = rns[j], rns[i]
	}
	return strings.ToLower(string(rns))
}
