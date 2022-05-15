package accounts

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

func (r *Repository) Create(publicId, fullName, position string) (*Account, error) {
	return nil, nil
}

func (r *Repository) ListWorkers() ([]Account, error) {
	var result []Account

	if err := r.db.Select(&result, "Select * from accounts where position = $1", WorkerRole); err != nil {
		return nil, err
	}
	return result, nil
}
