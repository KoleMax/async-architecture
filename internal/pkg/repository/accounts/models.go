package accounts

const WorkerRole = "worker"

type Account struct {
	Id       int    `db:"id"`
	PublicId string `db:"public_id"`
	Fullname string `db:"full_name"`
	Position string `db:"position"`
}
