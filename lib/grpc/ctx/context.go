package grpc_ctx

import (
	"context"

	"cp/lib/grpc/ptid"
	"google.golang.org/grpc/metadata"
)

func NewReqContext(ctx ...context.Context) context.Context {
	md := metadata.MD{}

	var ptid string
	if ctx != nil {
		if len(ptid) == 0 {
			ptid = grpc_ptid.GetPTIDFromContext(ctx[0])
		}
	}

	var userUuid string
	if ctx != nil {
		if len(userUuid) == 0 {
			userUuid = GetUserUuidFromContext(ctx[0])
		}
	}

	if len(ptid) == 0 {
		ptid = grpc_ptid.NewPTID()
	}

	md.Set(grpc_ptid.PassThroughIdField, ptid)
	md.Set(UserUuid, userUuid)

	if ctx != nil && len(ctx) != 0 {
		return metadata.NewOutgoingContext(ctx[0], md)
	}

	return metadata.NewOutgoingContext(context.Background(), md)
}
