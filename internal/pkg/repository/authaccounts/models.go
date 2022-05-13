package authaccounts

type AuthAccountCreateRow struct {
	Email    string `db:"email"`
	Password string `db:"password"`
	Fullname string `db:"full_name"`
	Position string `db:"position"`
}

type AuthAccountGetRow struct {
	Id       int    `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
	PublicId string `db:"public_id"`
	Fullname string `db:"full_name"`
	Position string `db:"position"`
}
