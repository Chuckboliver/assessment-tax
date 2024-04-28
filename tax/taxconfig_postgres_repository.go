package tax

import (
	"context"

	"github.com/jmoiron/sqlx"
)

var _ TaxConfigRepository = (*taxConfigPostgresRepository)(nil)

type taxConfigPostgresRepository struct {
	db sqlx.ExtContext
}

func NewTaxConfigPostgresRepository(db sqlx.ExtContext) TaxConfigRepository {
	return &taxConfigPostgresRepository{
		db: db,
	}
}

func (t *taxConfigPostgresRepository) FindByName(ctx context.Context, name string) (*Config, error) {
	sql := `
		SELECT name, value
		FROM tax_config
		WHERE name = $1
	`

	row := t.db.QueryRowxContext(ctx, sql, name)

	var config Config
	if err := row.Scan(&config.Name, &config.Value); err != nil {
		return nil, err
	}

	return &config, nil
}
