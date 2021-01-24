package endpoints

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/baxiang/soldiers-sortie/goblueprints-vault/services"
	"github.com/go-kit/kit/endpoint"
	"net/http"
)
type VaultEndpoints struct {
	HashEndpoint     endpoint.Endpoint
	ValidateEndpoint endpoint.Endpoint
}

type HashRequest struct {
	Password string `json:"password"`
}
type HashResponse struct {
	Hash string `json:"hash"`
	Err  string `json:"err,omitempty"`
}
type ValidateRequest struct {
	Password string `json:"password"`
	Hash     string `json:"hash"`
}
type ValidateResponse struct {
	Valid bool   `json:"valid"`
	Err   string `json:"err,omitempty"`
}

// Hash uses the HashEndpoint to hash a password.
func (e VaultEndpoints) Hash(ctx context.Context, password string) (string, error) {
	req := HashRequest{Password: password}
	resp, err := e.HashEndpoint(ctx, req)
	if err != nil {
		return "", err
	}
	hashResp := resp.(HashResponse)
	if hashResp.Err != "" {
		return "", errors.New(hashResp.Err)
	}
	return hashResp.Hash, nil
}

// Validate uses the ValidateEndpoint to validate a password and hash pair.
func (e VaultEndpoints) Validate(ctx context.Context, password,
	hash string) (bool, error) {
	req := ValidateRequest{Password: password, Hash: hash}
	resp, err := e.ValidateEndpoint(ctx, req)
	if err != nil {
		return false, err
	}
	validateResp := resp.(ValidateResponse)
	if validateResp.Err != "" {
		return false, errors.New(validateResp.Err)
	}
	return validateResp.Valid, nil
}

func MakeHashEndpoint(srv services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(HashRequest)
		v, err := srv.Hash(ctx, req.Password)
		if err != nil {
			return HashResponse{v, err.Error()}, nil
		}
		return HashResponse{v, ""}, nil
	}
}

func MakeValidateEndpoint(srv services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ValidateRequest)
		v, err := srv.Validate(ctx, req.Password, req.Hash)
		if err != nil {
			return ValidateResponse{false, err.Error()}, nil
		}
		return ValidateResponse{v, ""}, nil
	}
}

func DecodeHashRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req HashRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func DecodeValidateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req ValidateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}