package redis

// 定义redis key

const (
	// KeyPrefix redis key 注意使用命令空间的方式，方便查询和拆分
	KeyPrefix              = "bluebell:"
	KeyPostTimeZSet        = "post:time"   // zset;帖子及发帖时间
	KeyPostScoreZSet       = "post:score"  // zset;帖子及投票的分数
	keyPostVotedZSetPrefix = "post:voted:" // zset;记录用户及投票类型;参数是post_id
)

func getRedisKey(key string) string {
	return KeyPrefix + key
}
