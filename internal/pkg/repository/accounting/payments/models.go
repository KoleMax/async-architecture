package payments

import (
	"time"
)

type Payment struct {
	Id             int       `db:"id"`
	BillingCycleId int       `db:"billing_cycle_id"`
	AccountId      int       `db:"account_id"`
	Amount         int       `db:"amount"`
	Status         string    `db:"status"`
	CreatedAt      time.Time `db:"created_at"`
	EventCreatedAt time.Time `db:"event_created_at"`
}
