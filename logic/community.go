package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

// GetCommunityList 处理获取community业务逻辑
func GetCommunityList() (CommunityList []*models.Community, err error) {
	// 查数据库 查到所有community并返回
	return mysql.GetCommunityList()
}

// GetCommunityDetail 处理获取单个community详情业务逻辑
func GetCommunityDetail(idInt int) (CommunityDetail *models.CommunityDetail, err error) {
	// 查数据库 查到所有community并返回
	return mysql.GetCommunityDetailByID(idInt)
}
