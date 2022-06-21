package grpc_ptid

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

const PassThroughIdField = "ptid"

func NewPTID() string {
	return uuid.New().String()
}

func GetPTIDFromContext(ctx context.Context) string {

	if ctx == nil {
		return ""
	}

	ptid, ok := ctx.Value(PassThroughIdField).(string)
	if ok {
		return ptid
	}

	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		return md.Get(PassThroughIdField)[0]
	}

	md, ok = metadata.FromIncomingContext(ctx)
	if ok {
		return md.Get(PassThroughIdField)[0]
	}

	return ""

}
