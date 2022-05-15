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

func (r *Repository) Create(publicId, fullName, email string) error {
	if _, err := r.db.Exec("insert into accounts (public_id, full_name, email) values ($1, $2, $3)", publicId, fullName, email); err != nil {
		return err
	}
	return nil
}

func (r *Repository) SetBalance(id, balance int) error {
	if _, err := r.db.Exec("update  accounts set balance=$1 where id=$2", balance, id); err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetByPublicId(publicId string) (*Account, error) {
	var result Account

	if err := r.db.Get(&result, "select id, public_id, full_name, email, balance from accounts where email = $1", publicId); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}
