package handlersV1

import (
	"context"
	"go_http_barko/v1/domains"
)

type svcV1I interface {
	GetAllUserSvc(ctx context.Context) ([]domains.UserInfoResponseData, error)
}

type handlersV1 struct {
	svc svcV1I
}

func New(svcV1 svcV1I) *handlersV1 {
	return &handlersV1{
		svc: svcV1,
	}
}
