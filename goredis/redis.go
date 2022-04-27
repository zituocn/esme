/*
redis.go
go-redis v8 封装
sam
2022-04-25
*/

package goredis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/zituocn/esme/logx"
)

var (
	dbs           map[string]*redis.Client
	defaultDBName string
	ctx           = context.Background()
)

// RedisConfig redis 配置
type RedisConfig struct {
	Name     string // 连接名
	DB       int    // redis db num
	Host     string // 主机
	Port     int    // 端口
	Username string // 用户名，如果有
	Password string // 密码
	Pool     int    // 连接池大小
}

// InitDefaultDB 根据 rc 配置，初始化一个redis连接
func InitDefaultDB(rc *RedisConfig) (err error) {
	if rc == nil {
		err = errors.New("[redis] 没有需要init的配置")
		return
	}

	defaultDBName = rc.Name
	dbs = make(map[string]*redis.Client, 1)
	newRedis(rc)
	return
}

// InitDB 根据 rcs 配置，初始化多个redis连接
func InitDB(rcs []*RedisConfig) (err error) {
	if len(rcs) == 0 {
		err = errors.New("[redis] 没有需要init的配置")
		return
	}
	dbs = make(map[string]*redis.Client, len(rcs))
	for _, item := range rcs {
		newRedis(item)
	}

	return
}

// GetRDB 取默认的redis连接
//	只有一个redis连接时使用
func GetRDB() *redis.Client {
	rdb, ok := dbs[defaultDBName]
	if !ok {
		logx.Panicf("[redis] 未init，请阅读使用说明: https://github.com/zituocn/esme/goredis/README.md")
	}
	return rdb
}

// GetRDBByName 根据name取一个redis的连接
func GetRDBByName(name string) *redis.Client {
	m, ok := dbs[name]
	if !ok {
		logx.Panicf("[redis] 未init，请阅读使用说明: https://github.com/zituocn/esme/goredis/README.md")
	}
	return m
}

/*
private
*/

func (r *RedisConfig) string() string {
	return fmt.Sprintf("redis://%s:%s@%s:%d/%d", r.Name, r.Password, r.Host, r.Port, r.DB)
}

func newRedis(rc *RedisConfig) {
	var (
		rdb *redis.Client
	)
	if rc.Host == "" || rc.Port == 0 || rc.Name == "" {
		logx.Panicf("[redis]-[%s] 配置信息获取失败", rc.Name)
		return
	}
	if rc.DB < 0 {
		rc.DB = 0
	}
	if rc.Pool < 0 {
		rc.Pool = 10
	}
	opt := &redis.Options{
		Addr:         fmt.Sprintf("%s:%d", rc.Host, rc.Port),
		Username:     rc.Username,
		Password:     rc.Password,
		DB:           rc.DB,
		PoolSize:     rc.Pool,
		IdleTimeout:  30 * time.Second,
		DialTimeout:  5 * time.Second,
		MaxRetries:   -1,
		MinIdleConns: 10,
	}

	rdb = redis.NewClient(opt)

	// COMMAND ping
	for _, err := rdb.Ping(ctx).Result(); err != nil; {
		logx.Errorf("[redis]-%s 连接异常: %v", rc.string(), err)
		time.Sleep(5 * time.Second)
	}

	dbs[rc.Name] = rdb
}
