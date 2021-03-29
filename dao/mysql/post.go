package mysql

import (
	"bluebell/models"
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// CreatePost 向数据库保存用户创建的帖子
func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into 
		post(post_id, title, content, author_id, community_id)
		values(?, ?, ?, ?, ?)
		`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

// GetPostDetailByID 通过帖子ID获取帖子详情
func GetPostDetailByID(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
	from post
	where post_id=?
	`
	if err = db.Get(post, sqlStr, pid); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("查不到记录", zap.Error(err))
			err = ErrorInvalidID
		}
	}
	return
}

// GetPostList 获取帖子列表
func GetPostList(size, page int64) (postList []*models.Post, err error) {
	postList = make([]*models.Post, 0, 2)
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
	from post
	limit ?,?
	`
	err = db.Select(&postList, sqlStr, (page-1)*size, size)
	return
}

// 根据给定的id帖子查询帖子数据
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
	from post
	where post_id in (?)
	order by FIND_IN_SET(post_id, ?)
	`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	err = db.Select(&postList, query, args...)
	return
}
