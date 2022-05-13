package tasks

type Task struct {
	Id          int    `json:"id"`
	AssigneId   int    `json:"assigne_id"`
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

// Kafka

type BaseKafkaMessage struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

const (
	accountCreatedMsgType = "account.created"
)

type AccountCreatedMessage struct {
	PublicId string `json:"public_id"`
	Email    string `json:"email"`
	Fullname string `json:"full_name"`
	Position string `json:"position"`
}

const (
	taskAddedType     = "task.added"
	taskCompletedType = "task.completed"
	taskAssignedType  = "task.assigned"
)

type TaskAddedMessage struct {
	Id              int    `json:"id"`
	AssignePublicId string `json:"assigne_public_id"`
	Description     string `json:"descirption"`
}

type TaskAssignedMessage struct {
	Id              int    `json:"id"`
	AssignePublicId string `json:"assigne_public_id"`
}

type TaskCompletedMessage struct {
	Id              int    `json:"id"`
	AssignePublicId string `json:"assigne_public_id"`
	When            string `json:"when"`
}
