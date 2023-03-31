package cache

import (
	"context"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	DefaultMaxIdle   = 4
	DefaultMaxActive = 100
	DefaultPoolKey   = "default"
)

var (
	globalMapPool  = make(map[string]*redis.Pool)
	defaultOptions = Options{
		maxIdle:   DefaultMaxIdle,
		maxActive: DefaultMaxActive,
	}
)

func GetGlobalRedisPool() *redis.Pool {
	return GetRedisPoolByName(DefaultPoolKey)
}
func GetRedisPoolByName(name string) *redis.Pool {
	return globalMapPool[name]
}

type Options struct {
	maxIdle, maxActive int
}
type Option func(*Options)

func WithMaxIdle(maxIdle int) Option {
	return func(o *Options) {
		o.maxIdle = maxIdle
	}
}
func WithMaxActive(maxActive int) Option {
	return func(o *Options) {
		o.maxActive = maxActive
	}
}

type RedisConf struct {
	Name      string
	Server    string
	Pass      string
	MaxIdle   int
	MaxActive int
}

func NewPools(c []RedisConf) {
	for _, conf := range c {
		var options []Option
		if conf.MaxActive > 0 {
			options = append(options, WithMaxActive(conf.MaxActive))
		}
		if conf.MaxIdle > 0 {
			options = append(options, WithMaxIdle(conf.MaxIdle))
		}

		p, err := NewPool(conf.Server, conf.Pass, options...)
		if err != nil {
			panic(err)
		}
		nameKey := conf.Name
		if nameKey == "" {
			nameKey = DefaultPoolKey
		}
		globalMapPool[nameKey] = p
	}
}
func NewPool(server, pass string, options ...Option) (*redis.Pool, error) {
	dOptions := defaultOptions
	for _, o := range options {
		o(&dOptions)
	}
	p := &redis.Pool{
		MaxIdle:   dOptions.maxIdle,
		MaxActive: dOptions.maxActive,
		Dial: func() (redis.Conn, error) {
			diagOptions := make([]redis.DialOption, 0)
			if pass != "" {
				diagOptions = append(diagOptions, redis.DialPassword(pass))
			}
			rc, err := redis.Dial("tcp", server, diagOptions...)
			if err != nil {
				return nil, err
			}
			return rc, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		IdleTimeout: time.Second,
		Wait:        true,
	}
	if err := p.Get().Err(); err != nil {
		return nil, err
	}
	return p, nil
}
func StopAllPools() {
	for _, v := range globalMapPool {
		_ = v.Close()
	}
}

func HMSetContext(ctx context.Context, database int, key string, value map[interface{}]interface{}) (string, error) {
	conn, err := GetGlobalRedisPool().GetContext(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	err = conn.Send("SELECT", database)
	if err != nil {
		return "", err
	}
	return redis.String(conn.(redis.ConnWithContext).DoContext(ctx, "HMSET", redis.Args{}.Add(key).AddFlat(value)...))
}

func HMGetAllContext(ctx context.Context, database int, key string) (map[string]string, error) {
	conn, err := GetGlobalRedisPool().GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	err = conn.Send("SELECT", database)
	if err != nil {
		return nil, err
	}
	return redis.StringMap(conn.(redis.ConnWithContext).DoContext(ctx, "HGETALL", key))
}

func GetContext(ctx context.Context, database int, key string) (string, error) {
	conn, err := GetGlobalRedisPool().GetContext(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	err = conn.Send("SELECT", database)
	if err != nil {
		return "", err
	}
	return redis.String(conn.(redis.ConnWithContext).DoContext(ctx, "GET", key))
}

func MGetContext(ctx context.Context, database int, keys []interface{}) ([]string, error) {
	conn, err := GetGlobalRedisPool().GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	err = conn.Send("SELECT", database)
	if err != nil {
		return nil, err
	}
	return redis.Strings(conn.(redis.ConnWithContext).DoContext(ctx, "MGET", keys...))
}

func SetContext(ctx context.Context, database int, key string, value interface{}, expire uint32) (string, error) {
	conn, err := GetGlobalRedisPool().GetContext(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	err = conn.Send("SELECT", database)
	if err != nil {
		return "", err
	}
	args := []interface{}{
		key, value,
	}
	if expire > 0 {
		args = append(args, "EX")
		args = append(args, expire)
	}
	return redis.String(conn.(redis.ConnWithContext).DoContext(ctx, "SET", args...))
}
func MultiContext(ctx context.Context, database int, operations []string, args []redis.Args) ([]interface{}, error) {
	conn, err := GetGlobalRedisPool().GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	err = conn.Send("SELECT", database)
	if err != nil {
		return nil, err
	}
	err = conn.Send("MULTI")
	if err != nil {
		return nil, err
	}
	for k, arg := range args {
		err = conn.Send(operations[k], arg...)
	}
	return redis.Values(conn.(redis.ConnWithContext).DoContext(ctx, "EXEC"))
}

func ZAddContext(ctx context.Context, key string, value map[interface{}]interface{}) (int, error) {
	conn, err := GetGlobalRedisPool().GetContext(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	return redis.Int(conn.(redis.ConnWithContext).DoContext(ctx, "ZADD", redis.Args{}.Add(key).AddFlat(value)...))
}
func ZRangeByScoreContext(ctx context.Context, key string, min, max interface{}) ([]string, error) {
	conn, err := GetGlobalRedisPool().GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return redis.Strings(conn.(redis.ConnWithContext).DoContext(ctx, "ZRANGEBYSCORE", redis.Args{}.Add(key).Add(min).Add(max)...))
}
func ZRangeByScoreLimitContext(ctx context.Context, key string, min, max interface{}, offset, num int) ([]string, error) {
	conn, err := GetGlobalRedisPool().GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	args := redis.Args{}.Add(key).Add(min).Add(max).Add("LIMIT").Add(offset).Add(num)
	return redis.Strings(conn.(redis.ConnWithContext).DoContext(ctx, "ZRANGEBYSCORE", args...))
}
func ZRemContext(ctx context.Context, key, mem string) (int, error) {
	conn, err := GetGlobalRedisPool().GetContext(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	return redis.Int(conn.(redis.ConnWithContext).DoContext(ctx, "ZREM", redis.Args{}.Add(key).Add(mem)...))
}

func HSetContext(ctx context.Context, key, field string, value interface{}) (int, error) {
	conn, err := GetGlobalRedisPool().GetContext(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	return redis.Int(conn.(redis.ConnWithContext).DoContext(ctx, "HSET", redis.Args{}.Add(key).Add(field).Add(value)...))
}
func HGetContext(ctx context.Context, key, field string) (string, error) {
	conn, err := GetGlobalRedisPool().GetContext(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	return redis.String(conn.(redis.ConnWithContext).DoContext(ctx, "HGET", redis.Args{}.Add(key).Add(field)...))
}

func HDelContext(ctx context.Context, key, field string) (int, error) {
	conn, err := GetGlobalRedisPool().GetContext(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	return redis.Int(conn.(redis.ConnWithContext).DoContext(ctx, "HDEL", redis.Args{}.Add(key).Add(field)...))
}
func SetNXExpireContext(ctx context.Context, key string, value interface{}, expire time.Duration) (string, error) {
	conn, err := GetGlobalRedisPool().GetContext(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	args := []interface{}{
		key, value, "NX", "EX", uint32(expire / time.Second),
	}
	return redis.String(conn.(redis.ConnWithContext).DoContext(ctx, "SET", args...))
}
func SetNXContext(ctx context.Context, key string, value interface{}) (string, error) {
	conn, err := GetGlobalRedisPool().GetContext(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	args := []interface{}{
		key, value, "NX",
	}
	return redis.String(conn.(redis.ConnWithContext).DoContext(ctx, "SET", args...))
}
func GetNoDataBaseContext(ctx context.Context, key string) (string, error) {
	conn, err := GetGlobalRedisPool().GetContext(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	return redis.String(conn.(redis.ConnWithContext).DoContext(ctx, "GET", key))
}

func DelContext(ctx context.Context, key string) (int, error) {
	conn, err := GetGlobalRedisPool().GetContext(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	return redis.Int(conn.(redis.ConnWithContext).DoContext(ctx, "DEL", key))
}
