package handlersV1

import (
	"context"
)

type svcV1I interface {
	GetHelloWorldSvc(ctx context.Context) string
}

type handlersV1 struct {
	ctx context.Context
	svc svcV1I
}

func New(ctx context.Context, svcV1 svcV1I) *handlersV1 {
	return &handlersV1{
		ctx: ctx,
		svc: svcV1,
	}
}
