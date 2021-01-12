package job

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/hack-fan/skadi/types"
)

type Service struct {
	ctx    context.Context
	log    *zap.SugaredLogger
	kv     *redis.Client
	db     *gorm.DB
	rest   *resty.Client
	notify types.NotifyFunc
}

// NewService create a job service instance
func NewService() *Service {
	var s = &Service{
		ctx: context.Background(),
	}
	return s
}
