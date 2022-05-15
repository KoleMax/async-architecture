package tasks

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

func (r *Repository) Create(publicId string, costDone, costAssigne int, title, jira_id, description string) error {
	if _, err := r.db.Exec("insert into tasks (public_id, cost_done, cost_assigne, title, jira_id, description) values ($1, $2, $3, $4, $5, $6)", publicId, costDone, costAssigne, title, jira_id, description); err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetById(taskId int) (*Task, error) {
	var result Task

	if err := r.db.Get(&result, "select id, public_id, description, cost_done, cost_assigne from tasks where id = $1", taskId); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}
	return &result, nil
}

func (r *Repository) GetByPublicId(taskPublicId string) (*Task, error) {
	var result Task

	if err := r.db.Get(&result, "select id, public_id, description, cost_done, cost_assigne from tasks where public_id = $1", taskPublicId); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}
