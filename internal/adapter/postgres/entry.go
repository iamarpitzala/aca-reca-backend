package postgres

import "github.com/jmoiron/sqlx"

type entryRepo struct {
	db *sqlx.DB
}

func NewEntryRepo(db *sqlx.DB) *entryRepo {
	return &entryRepo{db: db}
}
