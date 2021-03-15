package event

import (
	"github.com/hack-fan/x/rdb"
	"go.uber.org/zap"

	"github.com/hack-fan/skadi/types"
)

const (
	CenterTypeRedis = "redis"
)

type Config struct {
	Type  string
	Redis rdb.Config
}

func NewEventCenter(log *zap.SugaredLogger, config Config) types.EventCenter {
	switch config.Type {
	case CenterTypeRedis:
		kv := rdb.New(config.Redis)
		return NewRedisEventCenter(kv, "skadi:event", log)
	default:
		return NewDefaultEventCenter(log)
	}
}
