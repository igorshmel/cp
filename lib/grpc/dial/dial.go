package grpc_dial

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/stats"
)

func NewConnection(ctx context.Context, target, tlsSrvOverride string, handler stats.Handler) (*grpc.ClientConn, error) {
	ctx, _ = context.WithTimeout(ctx, 5*time.Second)

	secure := grpc.WithInsecure()
	if len(tlsSrvOverride) != 0 {
		// TODO Подтянуть сертификаты
	}

	conn, err := grpc.DialContext(ctx, target, secure,
		grpc.WithBlock(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                5 * time.Minute,
			PermitWithoutStream: true,
		}),
		grpc.WithStatsHandler(handler))

	switch err {
	case nil:
	case context.DeadlineExceeded:
		return nil, fmt.Errorf("unable connect to gRPC service: %s", err.Error())
	default:
		return nil, err
	}

	return conn, nil
}
