package servicesV1

import (
	"context"
	"go_http_barko/utility/logger"
	"go_http_barko/v1/domains"
)

func (s *svcV1) GetAllUserSvc(ctx context.Context) ([]domains.UserInfoResponseData, error) {
	userInfoRepoResp, err := s.repoV1I.GetAllUserRepo(ctx)
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	userInfoResp := make([]domains.UserInfoResponseData, len(userInfoRepoResp))
	for i := range userInfoRepoResp {
		userInfoResp[i].Id = userInfoRepoResp[i].Id
		userInfoResp[i].FirstName = userInfoRepoResp[i].FirstName
		userInfoResp[i].LastName = userInfoRepoResp[i].LastName
		userInfoResp[i].CreatedDatetime = userInfoRepoResp[i].CreatedDatetime
		userInfoResp[i].UpdateDatetime = userInfoRepoResp[i].UpdateDatetime
	}

	return userInfoResp, nil
}
