package accounting

import (
	"time"
)

// Auth

const (
	AdminPosition      = "admin"
	ManagerPosition    = "manager"
	AccountantPosition = "accountant"
	WorkerPosition     = "worker"
)

type AuthAccount struct {
	PublicId string `json:"public_id"`
	Email    string `json:"email"`
	Fullname string `json:"full_name"`
	Position string `json:"position"`
}

// Api

const (
	TransactionTypeCredit   = "credit"
	TransactionTypeWriteOff = "write-off"
)

type Transaction struct {
	Id              int       `json:"id"`
	TaskId          int       `json:"task_id"`
	TaskDescription string    `json:"task_descriptioin"`
	EventCreatedAt  time.Time `json:"event_created_at"`
	Type            string    `json:"type"`
	Cost            int       `json:"cost"`
}

type Payment struct {
	Id             int       `db:"id"`
	Amount         int       `db:"amount"`
	Status         string    `db:"status"`
	CreatedAt      time.Time `db:"created_at"`
	EventCreatedAt time.Time `db:"event_created_at"`
}
