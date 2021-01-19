package kit

import "context"

func DecodeRequest(ctx context.Context, request interface{}) (interface{}, error) {
	return request, nil
}

func EncodeResponse(ctx context.Context, response interface{}) (interface{}, error) {
	return response, nil
}
