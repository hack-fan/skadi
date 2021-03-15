package event

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack/v5"
	"go.uber.org/zap"

	"github.com/hack-fan/skadi/types"
)

type RedisEventCenter struct {
	ctx context.Context // just blank context used for go-redis api
	kv  *redis.Client
	key string
	log *zap.SugaredLogger
}

func NewRedisEventCenter(kv *redis.Client, key string, log *zap.SugaredLogger) *RedisEventCenter {
	return &RedisEventCenter{
		ctx: context.Background(),
		kv:  kv,
		key: key,
		log: log,
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

func (ec *RedisEventCenter) StartWorker(ctx context.Context, handler types.EventHandler) {
	go ec.startWorker(ctx, handler)
}

func (ec *RedisEventCenter) startWorker(ctx context.Context, handler types.EventHandler) {
	ec.log.Debugf("event worker start at redis key: %s", ec.key)
LOOP:
	for {
		select {
		case <-ctx.Done():
			ec.log.Infof("event worker graceful shutdown at redis key: %s", ec.key)
			break LOOP
		default:
			event, err := ec.Get()
			if err != nil {
				ec.log.Errorf("error get event at %s: %s", ec.key, err)
				// perhaps redis down, wait for it
				time.Sleep(time.Minute)
				continue
			}
			if event == nil {
				// no new event, wait a moment
				time.Sleep(time.Second * 3)
				continue
			}
			err = handler(event)
			if err != nil {
				ec.log.Error(err)
				continue
			}
		}
	}
}
