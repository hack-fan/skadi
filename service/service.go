package service

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"github.com/hack-fan/skadi/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	ctx  context.Context
	kv   *redis.Client
	db   *gorm.DB
	rest *resty.Client
	log  *zap.SugaredLogger
	ev   types.EventCenter
}

// New create a job service instance
func New(kv *redis.Client, db *gorm.DB, rest *resty.Client, log *zap.SugaredLogger, ev types.EventCenter) *Service {
	var s = &Service{
		ctx:  context.Background(),
		kv:   kv,
		db:   db,
		rest: rest,
		log:  log,
		ev:   ev,
	}
	return s
}
