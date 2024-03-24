package repositoriesV1

import (
	"context"
	"fmt"
	"go_http_barko/utility/logger"
	"go_http_barko/v1/models"
)

func (r *repoV1) GetAllUserRepo(ctx context.Context) ([]models.UserInfo, error) {
	sqlRow, err := r.db.QueryContext(ctx, fmt.Sprintf("SELECT * FROM %v", models.USER_INFO_TABLE))
	if err != nil {
		logger.Error(ctx, err)
	}
	defer sqlRow.Close()

	var userInfo []models.UserInfo
	for sqlRow.Next() {
		var u models.UserInfo
		if err := sqlRow.Scan(&u.Id, &u.FirstName, &u.LastName, &u.CreatedDatetime, &u.UpdateDatetime); err != nil {
			return userInfo, err
		}
		userInfo = append(userInfo, u)
	}

	return userInfo, nil
}
