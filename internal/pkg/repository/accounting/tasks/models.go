package tasks

type Task struct {
	Id          int    `db:"id"`
	PublicId    string `db:"public_id"`
	Description string `db:"description"`
	CostDone    int    `db:"cost_done"`
	CostAssigne int    `db:"cost_assigne"`
}
