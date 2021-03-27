package mysql

import (
	"bluebell/models"
	"database/sql"

	"go.uber.org/zap"
)

// GetCommunityList 查询所有Community信息
func GetCommunityList() (CommunityList []*models.Community, err error) {
	sqlStr := `select community_id,community_name from community`
	if err = db.Select(&CommunityList, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("查不到记录", zap.Error(err))
			err = nil
		}
	}
	return
}

// GetCommunityDetailByID 查询所有Community信息
func GetCommunityDetailByID(id int64) (CommunityDetail *models.CommunityDetail, err error) {
	CommunityDetail = new(models.CommunityDetail)
	sqlStr := `
		select community_id,community_name,introduction,create_time 
		from community 
		where community_id=?
		`
	if err = db.Get(CommunityDetail, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("查不到记录", zap.Error(err))
			err = ErrorInvalidID
		}
	}
	return
}
