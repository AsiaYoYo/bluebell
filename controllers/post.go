package controllers

import (
	"bluebell/logic"
	"bluebell/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// CreatePostHandler 处理创建post请求
func CreatePostHandler(c *gin.Context) {
	// 1. 获取参数
	var p = new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		// 获取validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			ResponseError(c, CodeInvalidParam)
			return
		}
		zap.L().Error("create post with invalid param", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 获取当前用户的ID
	userID, err := getCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	// 2. 业务处理
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return

	}
	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 处理获取帖子详情请求
func GetPostDetailHandler(c *gin.Context) {
	// 1. 从路径中获取参数
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail failed", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 2. 业务处理
	data, err := logic.GetPostDetail(pid)
	if err != nil {
		zap.L().Error("logic.GetPostDetail failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler 处理获取帖子列表的请求
func GetPostListHandler(c *gin.Context) {
	// 1.获取分页参数
	size, page := getPostInfo(c)
	// 1. 业务处理
	data, err := logic.GetPostList(size, page)
	if err != nil {
		zap.L().Error("logic.GetPostList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 2. 返回响应
	ResponseSuccess(c, data)
}

// 根据创建时间排序 或者 按照分数排序
// GetPostListHandler2 处理根据path中的query参数获取帖子列表的请求
func GetPostListHandler2(c *gin.Context) {
	// 1.GET请求参数
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 with invaild param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 1. 业务处理
	data, err := logic.GetPostList2(p)
	if err != nil {
		zap.L().Error("logic.GetPostList2 failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 2. 返回响应
	ResponseSuccess(c, data)
}
