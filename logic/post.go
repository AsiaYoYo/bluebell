package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"

	"go.uber.org/zap"
)

// CreatePost 创建帖子逻辑处理函数
func CreatePost(p *models.Post) (err error) {
	// 1. 生成post id
	p.ID = snowflake.GenID()
	// 2. 把post保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	// 3. 把post创建时间保存到redis
	err = redis.CreatePost(p.ID)
	return
}

// GetPostDetail 获取帖子详情逻辑处理函数
func GetPostDetail(pid int64) (postDetail *models.APIPostDetail, err error) {
	// 查询帖子详情
	post, err := mysql.GetPostDetailByID(pid)
	if err != nil {
		zap.L().Error("get post detail failed",
			zap.Int64("post_id", pid),
			zap.Error(err))
		return
	}
	// 根据post详情中的AuthorID查询username
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("get user by id failed",
			zap.Int64("user_id", post.AuthorID),
			zap.Error(err))
		return
	}
	// 根据post详情中的CommunityID查询community
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("get community detail by id failed",
			zap.Int64("community_id", post.CommunityID),
			zap.Error(err))
		return
	}
	postDetail = &models.APIPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	return
}

// GetPostList 获取帖子列表逻辑处理函数
func GetPostList(size, page int64) (data []*models.APIPostDetail, err error) {
	// 查询帖子列表
	postList, err := mysql.GetPostList(size, page)
	if err != nil {
		zap.L().Error("get post list failed",
			zap.Error(err))
		return
	}
	data = make([]*models.APIPostDetail, 0, len(postList))
	for _, post := range postList {
		// 根据post详情中的AuthorID查询username
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("get user by id failed",
				zap.Int64("user_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		// 根据post详情中的CommunityID查询community
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("get community detail by id failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postDetail := &models.APIPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

// GetPostList2 升级版获取帖子列表逻辑处理函数
func GetPostList2(p *models.ParamPostList) (data []*models.APIPostDetail, err error) {
	// 1. 根据key去redis查询id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder return 0 data")
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("ids", ids))
	// 2. 根据id列表查询帖子详情
	// 返回的数据还要按照给定的顺序返回
	postList, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("postList", postList))
	VoteData, err := redis.GetPostScoreByIDs(ids)
	if err != nil {
		return
	}
	// 3. 返将帖子的作者及社区信息查询出来填充到帖子中
	for idx, post := range postList {
		// 根据post详情中的AuthorID查询username
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("get user by id failed",
				zap.Int64("user_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		// 根据post详情中的CommunityID查询community
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("get community detail by id failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postDetail := &models.APIPostDetail{
			AuthorName:      user.Username,
			VoteNum:         VoteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}
