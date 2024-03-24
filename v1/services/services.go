package servicesV1

import (
	"context"
	"go_http_barko/v1/models"
)

type repoV1Interface interface {
	GetAllUserRepo(ctx context.Context) ([]models.UserInfo, error)
}

type svcV1 struct {
	repoV1I repoV1Interface
}

func New(repoV1I repoV1Interface) *svcV1 {
	return &svcV1{repoV1I}
}
