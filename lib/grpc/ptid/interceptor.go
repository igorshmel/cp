package grpc_ptid

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func Interceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	// Переменная под идентификатор
	var ptid string

	// Вытаскиваем metadata (тут должен быть заголовок со сквозным идентификатором)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Unable to retrieve request metadata")
	}
	// Получаем заголовок с id
	ptidHeader := md.Get(PassThroughIdField)
	if len(ptidHeader) == 0 {
		// Если заголовок отсутствует - генерируем новый идентификатор
		ptid = uuid.New().String()
	} else {
		// Если заголовок есть - забираем идентификатор оттуда
		ptid = ptidHeader[0]
	}

	ctxWithPtid := context.WithValue(ctx, PassThroughIdField, ptid)

	return handler(ctxWithPtid, req)
}
