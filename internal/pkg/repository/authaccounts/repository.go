package authaccounts

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

func (r *Repository) Create(input *AuthAccountCreateRow) error {
	if _, err := r.db.NamedExec("insert into auth_accounts (email, password, full_name, position) values (:email, :password, :full_name, :position)", input); err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetByEmail(email string) (*AuthAccountGetRow, error) {
	result := new(AuthAccountGetRow)

	if err := r.db.Get(result, "select id, email, password, public_id, full_name, position from auth_accounts where email = $1", email); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}
