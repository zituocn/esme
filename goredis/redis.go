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

// RedisConfig redis config
type RedisConfig struct {
	Name     string // name
	DB       int    // redis db num
	Host     string // host
	Port     int    // port
	Username string // username
	Password string // password
	Pool     int    // pool
}

// InitDefaultDB Initialize a redis connection
func InitDefaultDB(rc *RedisConfig) (err error) {
	if rc == nil {
		err = errors.New("[redis] no configuration to initialize")
		return
	}

	defaultDBName = rc.Name
	dbs = make(map[string]*redis.Client, 1)
	newRedis(rc)
	return
}

// InitDB initialize multiple redis connections
func InitDB(rcs []*RedisConfig) (err error) {
	if len(rcs) == 0 {
		err = errors.New("[redis] no configuration to initialize")
		return
	}
	dbs = make(map[string]*redis.Client, len(rcs))
	for _, item := range rcs {
		newRedis(item)
	}

	return
}

// GetRDB get the default redis connection
func GetRDB() *redis.Client {
	rdb, ok := dbs[defaultDBName]
	if !ok {
		logx.Panicf("[redis] not initialized, please read the docs: https://github.com/zituocn/esme/goredis/README.md")
	}
	return rdb
}

// GetRDBByName get a redis connection by name
func GetRDBByName(name string) *redis.Client {
	m, ok := dbs[name]
	if !ok {
		logx.Panicf("[redis] not initialized, please read the docs: https://github.com/zituocn/esme/goredis/README.md")
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
		logx.Panicf("[redis]-[%s] failed to read configuration information", rc.Name)
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
		logx.Errorf("[redis]-%s connection exception: %s", rc.string(), err.Error())
		time.Sleep(5 * time.Second)
	}

	dbs[rc.Name] = rdb
}
