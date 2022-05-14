package tasks

const (
	TaskStatusActive    = "active"
	TaskStatusCompleted = "completed"
)

type Task struct {
	Id          int    `db:"id"`
	Description string `db:"description"`
	Status      string `db:"status"`
	AssigneId   int    `db:"assignee_id"`
	Version     int    `db:"version"`
}
