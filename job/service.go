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
	kv     *redis.Client
	db     *gorm.DB
	rest   *resty.Client
	log    *zap.SugaredLogger
	notify types.NotifyFunc
}

// NewService create a job service instance
func NewService(kv *redis.Client, db *gorm.DB, rest *resty.Client, log *zap.SugaredLogger) *Service {
	var s = &Service{
		ctx:  context.Background(),
		kv:   kv,
		db:   db,
		rest: rest,
		log:  log,
		notify: func(err error) {
			log.Error(err)
		},
	}
	return s
}

func (s *Service) SetNotifyFunc(notifyFunc types.NotifyFunc) {
	s.notify = notifyFunc
}
