package transactions

import (
	"time"

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

func (r *Repository) Create(accountId, taskId int, eventCreatedAt time.Time, transactionType string) error {
	_, err := r.db.Exec(`insert into transactions (account_id, task_id, event_created_at, type) values ($1, $2, $3, $4)`, accountId, taskId, eventCreatedAt, transactionType)
	return err
}

func (r *Repository) GetByAccountId(accountId int) ([]Transaction, error) {
	var result []Transaction

	if err := r.db.Select(result, `
select 
	tr.id id, 
	tr.account_id account_id, 
	tr.task_id task_id, 
	tr.event_created_at event_created_at, 
	tr.type type,
	t.description task_description,
	t.cost_done cost_done,
	t.cost_assigne cost_assigne
from 
	transactions tr 
		join
	tasks t
		on tr.task_id = t.id  
where 
	tr.account_id = $1
`, accountId); err != nil {
		return nil, err
	}
	return result, nil
}
