package models

import "time"

type Post struct {
	ID          int64     `json:"id,string" db:"post_id"`
	AuthorID    int64     `json:"author_id,string" db:"author_id"`
	CommunityID int64     `json:"community_id" db:"community_id" binding:"required"`
	Status      int32     `json:"status" db:"status"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
}

// APIPostDetail 定义帖子详情结构体
type APIPostDetail struct {
	AuthorName string `json:"author_name"`
	VoteNum    int64  `json:"vote_num"`
	*Post
	*CommunityDetail `json:"community"`
}
