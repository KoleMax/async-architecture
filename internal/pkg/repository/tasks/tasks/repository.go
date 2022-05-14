package tasks

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
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

func (r *Repository) Shuffle(assigneIds []int) error {
	return nil
}

func (r *Repository) Create(assigneId int, description string) (*Task, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("tasks").
		Columns(
			"assignee_id",
			"description",
		).Values(
		assigneId,
		description,
	).Suffix("returning id, description, assignee_id, status")
	stmt, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("builder.ToSql: %v", err)
	}

	var result Task
	if err = r.db.Get(&result, stmt, args...); err != nil {
		return nil, err
	}
	return &result, err
}

func (r *Repository) List(assignePublicId string) ([]Task, error) {
	var result []Task

	if err := r.db.Select(&result, "Select t.id, t.description, t.status, t.version, t.assignee_id from tasks t join accounts a on t.assignee_id = a.id where a.public_id = $1", assignePublicId); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repository) ListActive() ([]Task, error) {
	var result []Task

	if err := r.db.Select(&result, "Select id, description, status, version, assignee_id from tasks where status = $1", TaskStatusActive); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repository) ListAll() ([]Task, error) {
	var result []Task

	if err := r.db.Select(&result, "Select id, description, status, version, assignee_id from tasks"); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repository) ListAсtive() ([]Task, error) {
	var result []Task

	if err := r.db.Select(&result, fmt.Sprintf("Select id, description, status, version, assignee_id from tasks where status = $1", TaskStatusActive)); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repository) Complete(id int) error {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update("tasks").
		Set("status", TaskStatusCompleted).
		Where(sq.Eq{"id": id})

	stmt, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("builder.ToSql: %v", err)
	}

	_, err = r.db.Exec(stmt, args...)
	return err
}

func (r *Repository) Assigne(id, assigneId int) error {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update("tasks").
		Set("assignee_id", assigneId).
		Where(sq.Eq{"id": id})

	stmt, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("builder.ToSql: %v", err)
	}

	_, err = r.db.Exec(stmt, args...)
	return err
}
