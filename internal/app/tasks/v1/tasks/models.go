package tasks

type Task struct {
	Id          int    `json:"id"`
	AssigneId   int    `json:"assigne_id"`
	JiraId      string `json:"jira_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

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
