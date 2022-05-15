package transactions

import "time"

const (
	TransactionTypeDebit  = "debit"
	TransactionTypeCredit = "credit"
)

type Transaction struct {
	Id              int       `db:"id"`
	AccountId       int       `db:"account_id"`
	TaskId          int       `db:"task_id"`
	TaskDescription string    `db:"task_descriptioin"`
	EventCreatedAt  time.Time `db:"event_created_at"`
	Type            string    `db:"type"`
	CostAssigne     int       `db:"cost_assigne"`
	CostDone        int       `db:"cost_done"`
}
