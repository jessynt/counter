package counter

import (
	"context"

	"github.com/gomodule/redigo/redis"
)

type RedisCounter struct {
	pool *redis.Pool

	keyPrefix string
}

func NewRedisCounter(redisPool *redis.Pool, keyPrefix string) *RedisCounter {
	return &RedisCounter{
		pool:      redisPool,
		keyPrefix: keyPrefix,
	}
}

func (rc *RedisCounter) Incr(ctx context.Context, id string) error {
	conn, err := rc.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("INCR", rc.getKeyById(id))
	return err
}

func (rc *RedisCounter) Get(ctx context.Context, id string) (int64, error) {
	conn, err := rc.pool.GetContext(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	count, err := redis.Int64(conn.Do("GET", rc.getKeyById(id)))
	switch err {
	case nil:
		return count, nil
	case redis.ErrNil:
		return 0, nil
	default:
		return 0, err
	}
}

func (rc RedisCounter) getKeyById(id string) interface{} {
	return rc.keyPrefix + id
}
