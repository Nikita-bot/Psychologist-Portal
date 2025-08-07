package dal

import "database/sql"

type Repo struct {
	db *sql.DB
}

func InitRepo(db *sql.DB) *Repo {
	return &Repo{
		db: db,
	}
}
