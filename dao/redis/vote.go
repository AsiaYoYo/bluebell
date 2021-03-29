package redis

import (
	"errors"
	"math"
	"time"

	"github.com/go-redis/redis"
)

/* 投票的几种情况：
direction=1时，有两种情况：
	1. 之前没有投过票，现在投赞成票   --> 更新分数和投票记录 差值的绝对值：1 +432
	2. 之前投反对票，现在改投赞成票   --> 更新分数和投票记录 差值的绝对值：2 +432*2
direction=0时，有两种情况
	1. 之前投过反对票，现在要取消投票 --> 更新分数和投票记录 差值的绝对值：1 +432
	2. 之前投过赞成票，现在要取消投票 --> 更新分数和投票记录 差值的绝对值：1 -432
direction=-1时，有两种情况：
	1. 之前没有投过票，现在投反对票   --> 更新分数和投票记录 差值的绝对值：1 -432
	2. 之前投赞成票，现在改投反对票   --> 更新分数和投票记录 差值的绝对值：2 -432*2

投票的限制：
每个帖子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许再投票了。
	1. 到期之后将redis中保存的赞成票数存到mysql表中
	2. 到期之后删除那个 KeyPostVotedZSetPrefix
*/

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 // 每一票值多少分
)

var (
	ErrorVoteTimeExpire = errors.New("投票时间已过")
	ErrorVoteRepeated   = errors.New("不能重复投票")
)

// CreatePost 创建post时记录帖子原始分数
func CreatePost(postID int64) error {
	pipeline := rdb.TxPipeline()
	// 帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	// 帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	// 事务，一起执行
	_, err := pipeline.Exec()
	return err
}

// VoteForPost 对帖子进行投票操作存入redis
func VoteForPost(userID, postID string, value float64) error {
	// 1. 判断投票的限制
	// 去redis取帖子发布的时间
	postTime := rdb.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrorVoteTimeExpire
	}
	// 2. 更新分数
	// 先查当前用户给当前帖子的投票记录
	ov := rdb.ZScore(getRedisKey(keyPostVotedZSetPrefix+postID), userID).Val()
	if value == ov {
		return ErrorVoteRepeated
	}
	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value) // 计算两次投票的差值
	pipeline := rdb.TxPipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)
	// 3. 记录用户为该帖子投票的数据
	if value == 0 {
		pipeline.ZRem(getRedisKey(keyPostVotedZSetPrefix+postID), userID)
	} else {
		pipeline.ZAdd(getRedisKey(keyPostVotedZSetPrefix+postID), redis.Z{
			Score:  value, // 赞成票还是反对票
			Member: userID,
		})
	}
	_, err := pipeline.Exec()
	return err
}
