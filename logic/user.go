package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/snowflake"
)

// SignUp 用户注册逻辑处理
func SignUp(p *models.ParamSignUp) {
	// 1. 检查用户名是否重复
	mysql.QueryUserByUsername()
	// 2. 生成用户ID
	snowflake.GenID()
	// 2. 将数据写入数据库
	mysql.InsertUser()
}
