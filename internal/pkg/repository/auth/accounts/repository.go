package accounts

import (
	"database/sql"

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

func (r *Repository) Create(input *AccountCreateRow) error {
	if _, err := r.db.NamedExec("insert into accounts (email, password, full_name, position) values (:email, :password, :full_name, :position)", input); err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetByEmail(email string) (*AccountGetRow, error) {
	result := new(AccountGetRow)

	if err := r.db.Get(result, "select id, email, password, public_id, full_name, position from accounts where email = $1", email); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}
