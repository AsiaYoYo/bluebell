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
)

func CreatePost(postID int64) error {
	_, err := rdb.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	}).Result()
	return err
}

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
	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value) // 计算两次投票的差值
	_, err := rdb.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID).Result()
	if err != nil {
		return err
	}
	// 3. 记录用户为该帖子投票的数据
	if value == 0 {
		_, err = rdb.ZRem(getRedisKey(keyPostVotedZSetPrefix+postID), userID).Result()
	} else {
		_, err = rdb.ZAdd(getRedisKey(keyPostVotedZSetPrefix+postID), redis.Z{
			Score:  value, // 赞成票还是反对票
			Member: userID,
		}).Result()
	}
	return err
}
