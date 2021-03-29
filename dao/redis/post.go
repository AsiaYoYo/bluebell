package redis

import (
	"bluebell/models"

	"github.com/go-redis/redis"
)

// GetPostIDsInOrder 根据用户携带的order参数返回帖子id列表，分数从大到小
func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 从redis获取id
	// 根据用户请求中携带的order参数确定key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	// 2. 确认查询的索引起始点
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	// 3. ZREVRANGE 查询 按分数查询指定数量的元素
	return rdb.ZRevRange(key, start, end).Result()
}

// GetPostScoreByIDs 根据postid获取帖子分数
func GetPostScoreByIDs(ids []string) (data []int64, err error) {
	// 1. 循环所有key放入pipeline中
	pipeline := rdb.TxPipeline()
	for _, id := range ids {
		key := getRedisKey(keyPostVotedZSetPrefix + id)
		pipeline.ZCount(key, "1", "1")
	}
	// 2. pipeline执行
	cmders, err := pipeline.Exec()
	if err != nil {
		return
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}
