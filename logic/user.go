package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
)

// SignUp 用户注册逻辑处理
func SignUp(p *models.ParamSignUp) (err error) {
	// 1. 检查用户名是否重复
	if err = mysql.CheckUserExist(p.Username); err != nil {
		return
	}
	// 2. 生成用户ID
	UserID := snowflake.GenID()
	// 构造一个User实例
	user := &models.User{
		UserID:   UserID,
		Username: p.Username,
		Password: p.Password,
	}
	// 3. 将数据写入数据库
	return mysql.InsertUser(user)
}

// Login 用户登录逻辑处理
func Login(u *models.User) (userID int64, token string, err error) {
	// 构造一个Login实例
	user := &models.User{
		Username: u.Username,
		Password: u.Password,
	}
	// 1. 检查用户名和密码是否正确
	if err := mysql.Login(user); err != nil {
		return 0, "", err
	}
	// 2. 生成JWT
	token, err = jwt.GenToken(user.UserID, user.Username)
	return user.UserID, token, err
}
