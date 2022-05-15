package payments

import (
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(billingCycleId, accountId, amount int, status string) error {
	if _, err := r.db.Exec("insert into payments (billing_cycle_id, account_id, amount, status) values ($1, $2, $3, $4)", billingCycleId, accountId, amount, status); err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetByAccountId(accountId int) ([]Payment, error) {
	var result []Payment
	if err := r.db.Select(&result, "select id, billing_cycle_id, account_id, amount, status, created_at from payments p join billing_cycles where account_id = $1", accountId); err != nil {
		return nil, err
	}
	return result, nil
}
