package transport

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/baxiang/soldiers-sortie/string-server/endpoints"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
)

func GetIsPalHandler(ep endpoint.Endpoint)*httptransport.Server{
	return httptransport.NewServer(
		ep,
		decodeGetIsPalRequest,
		encodeGetIsPalResponse,
	)
}
func decodeGetIsPalRequest(_ context.Context,r *http.Request)(interface{},error){
	var req *endpoints.IsPalRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req,err
}

func encodeGetIsPalResponse(_ context.Context,w http.ResponseWriter,resp interface{})error{
	resp, ok := resp.(*endpoints.IsPalResponse)
	if !ok {
		return errors.New("error decoding")
	}
	return json.NewEncoder(w).Encode(resp)
}

func GetReverseHandler(ep endpoint.Endpoint) *httptransport.Server {
	return httptransport.NewServer(
		ep,
		decodeGetReverseRequest,
		encodeGetReverseResponse,
	)
}

func decodeGetReverseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req *endpoints.ReverseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func encodeGetReverseResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	resp, ok := response.(*endpoints.ReverseResponse)
	if !ok {
		return errors.New("error decoding")
	}
	return json.NewEncoder(w).Encode(resp)
}
