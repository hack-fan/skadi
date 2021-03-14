package event

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/hack-fan/skadi/types"
)

type RedisEventCenter struct {
	ctx context.Context // just blank context used for go-redis api
	kv  *redis.Client
	key string
}

func NewRedisEventCenter(kv *redis.Client, key string) *RedisEventCenter {
	return &RedisEventCenter{
		ctx: context.Background(),
		kv:  kv,
		key: key,
	}
}

func (ec *RedisEventCenter) Pub(e *types.Event) error {
	data, err := msgpack.Marshal(e)
	if err != nil {
		return fmt.Errorf("marshal event failed: %w", err)
	}
	err = ec.kv.LPush(ec.ctx, ec.key, data).Err()
	if err != nil {
		return fmt.Errorf("push event to redis list failed: %w", err)
	}
	return nil
}

func (ec *RedisEventCenter) Get() (*types.Event, error) {
	data, err := ec.kv.RPop(ec.ctx, ec.key).Bytes()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("get event from redis failed: %w", err)
	}
	var event = new(types.Event)
	err = msgpack.Unmarshal(data, event)
	if err != nil {
		return nil, fmt.Errorf("unmarshal event failed: %w", err)
	}
	return event, nil
}
