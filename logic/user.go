package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
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
