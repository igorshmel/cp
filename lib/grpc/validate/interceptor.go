package grpc_validate

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/*
Interceptor для валидации запроса, если у него описан метод Validate()
*/
func Interceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	v, ok := req.(ValidatableRequest)
	if ok {
		err := v.Validate()
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
	}

	return handler(ctx, req)
}
