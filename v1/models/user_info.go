package models

import "time"

const USER_INFO_TABLE = "user_info"

type UserInfo struct {
	Id              string
	FirstName       string
	LastName        string
	CreatedDatetime time.Time
	UpdateDatetime  time.Time
}
