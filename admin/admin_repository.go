package admin

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type adminRepository struct {
	db sqlx.ExtContext
}

var _ AdminRepository = (*adminRepository)(nil)

func NewAdminRepository(db sqlx.ExtContext) AdminRepository {
	return &adminRepository{
		db: db,
	}
}

func (r *adminRepository) UpdatePersonalDeduction(ctx context.Context, personalDeduction float64) error {
	sql := `
		UDPATE tax_config
		SET
			value = $1
		WHERE name = 'personal_deduction'
	`

	_, err := r.db.ExecContext(ctx, sql, personalDeduction)
	return err
}
