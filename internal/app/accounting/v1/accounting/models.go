package accounting

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
