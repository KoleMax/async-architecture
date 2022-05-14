package auth

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

const (
	accountCreatedMsgType = "account.created"
)

type BaseKafkaMessage struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

type AccountCreatedMessage struct {
	PublicId string `json:"public_id"`
	Email    string `json:"email"`
	Fullname string `json:"full_name"`
	Position string `json:"position"`
}
