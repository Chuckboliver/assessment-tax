package postgres

import (
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func New(dbURL string) (*sqlx.DB, error) {
	return sqlx.Open("postgres", dbURL)
}
