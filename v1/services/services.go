package servicesV1

import "context"

type repoV1Interface interface {
	GetHelloRepo(ctx context.Context, str string) string
}

type svcV1 struct {
	repoV1I repoV1Interface
}

func New(repoV1I repoV1Interface) *svcV1 {
	return &svcV1{repoV1I}
}
