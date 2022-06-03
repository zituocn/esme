/*
queue_redis.go
task queue implementation in redis
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

// RedisQueue task queue in redis
type RedisQueue struct {

	// redis list key
	key string

	// rdb redis client
	rdb *redis.Client
}

// NewRedisQueue use redis configuration
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

// Add add a task to the queue
func (q *RedisQueue) Add(task *Task) {
	b, err := json.Marshal(task)
	if err != nil {
		logx.Errorf("serialization task failed : %v", err)
		return
	}
	err = q.rdb.RPush(ctx, q.key, string(b)).Err()
	if err != nil {
		logx.Errorf("failed to add task to queue : %v", err)
		return
	}
}

// AddTasks add multiple tasks to the queue
func (q *RedisQueue) AddTasks(list []*Task) {
	for _, item := range list {
		q.Add(item)
	}
}

// Pop get a task while removing it from the queue
//	from left
func (q *RedisQueue) Pop() *Task {
	s, err := q.rdb.LPop(ctx, q.key).Result()
	if err != nil && err != redis.Nil {
		logx.Errorf("pop failed : %v", err)
		return nil
	}
	if len(s) == 0 {
		return nil
	}
	task := new(Task)

	err = json.Unmarshal([]byte(s), &task)
	if err != nil {
		logx.Errorf("return serialization task failure: %s", err.Error())
		return nil
	}
	return task
}

// Clear clear all tasks
func (q *RedisQueue) Clear() bool {
	i, err := q.rdb.Del(ctx, q.key).Result()
	if err != nil {
		logx.Errorf("Clear: %s", err.Error())
		return false
	}
	if i > 0 {
		return true
	}
	return false
}

// IsEmpty returns whether the queue is empty
func (q *RedisQueue) IsEmpty() bool {
	if q.Size() == 0 {
		return true
	}
	return false
}

// Size returns queue length
func (q *RedisQueue) Size() int {
	i, err := q.rdb.LLen(ctx, q.key).Result()
	if err != nil {
		return 0
	}
	return int(i)
}

func (q *RedisQueue) Print() {
}
