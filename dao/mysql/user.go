package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
)

const secret = "asiayoyo"

// CheckUserExist 用过用户名查询用户
func CheckUserExist(username string) (err error) {
	sqlStr := `select user_id from user where username=?`
	var count int
	err = db.Get(&count, sqlStr, username)
	if err == sql.ErrNoRows && count == 0 {
		return nil
	}
	if count > 0 {
		err = errors.New("用户已存在")
		return
	}
	return
}

// InsertUser 插入用户数据
func InsertUser(user *models.User) (err error) {
	user.Password = encryptPassword(user.Password)
	sqlStr := `insert into user(user_id,username,password) values(?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return
}

func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}
