package models

type User struct {
	UserID   int64  `db:"user_id" json:"user_id,string"`
	Username string `db:"username" json:"username" binding:"required"`
	Password string `db:"password" json:"password" binding:"required"`
}
