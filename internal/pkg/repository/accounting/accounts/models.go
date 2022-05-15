package accounts

type Account struct {
	Id       int    `db:"id"`
	PublicId string `db:"public_id"`
	FullName string `db:"full_name"`
	Email    string `db:"email"`
	Balance  int    `db:"balance"`
}
