/*
queue_redis.go
redis中的任务队列实现
sam
2022-04-25
*/

package esme

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/zituocn/esme/goredis"
	"github.com/zituocn/esme/logx"
)

var (
	ctx = context.Background()
)

// RedisQueue redis中的任务队列
type RedisQueue struct {

	// redis list key
	key string

	// rdb redis client
	rdb *redis.Client
}

// NewRedisQueue 使用redis配置
func NewRedisQueue(key string, rc *goredis.RedisConfig) TodoQueue {
	err := goredis.InitDefaultDB(rc)
	if err != nil {
		logx.Error(err)
		return nil
	}

	return &RedisQueue{
		key: key,
		rdb: goredis.GetRDB(),
	}
}

// Add 添加一个任务到队列中
func (q *RedisQueue) Add(task *Task) {
	b, err := json.Marshal(task)
	if err != nil {
		logx.Errorf("序列化任务失败 : %v", err)
		return
	}
	err = q.rdb.RPush(ctx, q.key, string(b)).Err()
	if err != nil {
		logx.Errorf("向队列添加任务失败 : %v", err)
		return
	}
}

// AddTasks 添加多个任务到队列中
func (q *RedisQueue) AddTasks(list []*Task) {
	for _, item := range list {
		q.Add(item)
	}
}

// Pop 到一个任务，同时从队列中移出它
//	从队列的左侧pop元素
func (q *RedisQueue) Pop() *Task {
	s, err := q.rdb.LPop(ctx, q.key).Result()
	if err != nil && err != redis.Nil {
		logx.Errorf("Pop失败 : %v", err)
		return nil
	}
	if len(s) == 0 {
		return nil
	}
	task := new(Task)

	err = json.Unmarshal([]byte(s), &task)
	if err != nil {
		logx.Errorf("返回序列化任务失败: %v", err)
		return nil
	}
	return task
}

// Clear 清理掉所有任务
func (q *RedisQueue) Clear() bool {
	i, err := q.rdb.Del(ctx, q.key).Result()
	if err != nil {
		logx.Errorf("Clear : %v", err)
		return false
	}
	if i > 0 {
		return true
	}
	return false
}

// IsEmpty 返回队列是否为空
func (q *RedisQueue) IsEmpty() bool {
	if q.Size() == 0 {
		return true
	}
	return false
}

// Size 返回队列的长度
func (q *RedisQueue) Size() int {
	i, err := q.rdb.LLen(ctx, q.key).Result()
	if err != nil {
		return 0
	}
	return int(i)
}

func (q *RedisQueue) Print() {
}
