package domains

import "time"

type UserInfoResponseData struct {
	Id              string    `json:"id"`
	FirstName       string    `json:"firstname"`
	LastName        string    `json:"lastname"`
	CreatedDatetime time.Time `json:"created_datetime"`
	UpdateDatetime  time.Time `json:"updated_datetime"`
}
