package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
)

const secret = "asiayoyo"

// CheckUserExist 根据用户名查询用户
func CheckUserExist(username string) (err error) {
	sqlStr := `select user_id from user where username=?`
	var count int
	err = db.Get(&count, sqlStr, username)
	if err == sql.ErrNoRows && count == 0 {
		return nil
	}
	if count > 0 {
		return ErrorUserExist
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

// Login 用户登录查询
func Login(user *models.User) (err error) {
	oPassword := user.Password
	// 1.查询用户的密码
	sqlStr := `select user_id,username,password from user where username=?`
	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil {
		// 查询数据库失败
		return
	}
	// 2.login Password和数据库中的密码进行比较
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return
}

// GetUserByID 根据用户ID查询用户信息
func GetUserByID(uid int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id,username from user where user_id=?`
	err = db.Get(user, sqlStr, uid)
	return
}
