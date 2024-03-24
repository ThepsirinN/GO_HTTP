package handlersV1

import (
	"context"
)

type svcV1I interface {
	GetHelloWorldSvc(ctx context.Context) string
}

type handlersV1 struct {
	svc svcV1I
}

func New(svcV1 svcV1I) *handlersV1 {
	return &handlersV1{
		svc: svcV1,
	}
}
