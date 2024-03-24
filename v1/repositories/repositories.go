package repositoriesV1

import "database/sql"

type repoV1 struct {
	db *sql.DB
}

func New(db *sql.DB) *repoV1 {
	return &repoV1{db}
}
