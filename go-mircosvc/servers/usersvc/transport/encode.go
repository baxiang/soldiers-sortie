package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	userPb "github.com/baxiang/soldiers-sortie/go-mircosvc/pb"
	"github.com/baxiang/soldiers-sortie/go-mircosvc/pkg/common"
	"github.com/baxiang/soldiers-sortie/go-mircosvc/servers/usersvc/endpoints"
	kitTransport "github.com/go-kit/kit/transport/http"
)

func encodeResponseSetCookie(_ context.Context, w http.ResponseWriter, response interface{}) error {
	if headerer, ok := response.(kitTransport.Headerer); ok {
		for k, values := range headerer.Headers() {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
	}
	// cookie := &http.Cookie{
	// 	Name:     common.AuthHeaderKey,
	// 	Value:    "test",
	// 	Path:     "/",
	// 	HttpOnly: true,
	// 	MaxAge:   int(common.MaxAge),
	// }
	// http.SetCookie(w, cookie)

	return json.NewEncoder(w).Encode(response)
}

// GetUser
func encodeGRPCGetUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	r, ok := request.(string)
	if !ok {
		return nil, errors.New("encodeGRPCGetUserRequest: interface conversion error")
	}

	return &userPb.GetUserRequest{Uid: r}, nil
}

// Login
func encodeGRPCLoginRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(endpoints.LoginRequest)
	if !ok {
		return nil, errors.New("encodeGRPCLoginRequest: interface conversion error")
	}
	return &userPb.LoginRequest{Username: req.Username, Password: req.Password}, nil
}

// SendCode
// ...

// Register
func encodeGRPCRegisterRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(endpoints.RegisterRequest)
	if !ok {
		return nil, errors.New("encodeGRPCRegisterRequest: interface conversion error")
	}
	return &userPb.RegisterRequest{Username: req.Username, Password: req.Password, CodeID: req.CodeID}, nil
}

// UserList
func encodeGRPCUserListRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(endpoints.UserListRequest)
	if !ok {
		return nil, errors.New("encodeGRPCUserListRequest: interface conversion error")
	}
	return &userPb.UserListRequest{Page: req.Page, Size: req.Size}, nil
}

func encodeGRPCUserListResponse(_ context.Context, response interface{}) (interface{}, error) {
	res, ok := response.(common.Response)
	if !ok {
		return nil, errors.New("encodeGRPCUserListResponse: interface conversion error")
	}

	data := res.Data.(*userPb.UserListResponse)

	// us := make([]*userPb.UserResponse, 0)
	// for _, v := range data.Data {
	// 	u := new(userPb.UserResponse)
	// 	if err := utils.StructCopy(v, u); err != nil {
	// 		return nil, err
	// 	}
	// 	us = append(us, u)
	// }

	return data, nil
}

// Logout
func encodeGRPCLogoutRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(endpoints.LogoutRequest)
	if !ok {
		return nil, errors.New("encodeGRPCLogoutRequest: interface conversion error")
	}
	return &userPb.LogoutRequest{Sid: req.SID}, nil
}