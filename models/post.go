package models

type Post struct {
	ID          int64  `json:"post_id" db:"post_id"`
	AuthorID    int64  `json:"author_id" db:"author_id"`
	CommunityID int    `json:"community_id" db:"community_id" binding:"required"`
	Title       string `json:"title" db:"title" binding:"required"`
	Content     string `json:"content" db:"content" binding:"required"`
}
