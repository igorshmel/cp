package grpc_ctx

import (
	"context"

	"google.golang.org/grpc/metadata"
)

const UserUuid = "useruuid"

func GetUserUuidFromContext(ctx context.Context) string {

	if ctx == nil {
		return ""
	}

	ptid, ok := ctx.Value(UserUuid).(string)
	if ok {
		return ptid
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		return md.Get(UserUuid)[0]
	}

	md, ok = metadata.FromOutgoingContext(ctx)
	if ok {
		return md.Get(UserUuid)[0]
	}

	return ""
}
