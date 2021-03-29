package service

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	ctx  context.Context
	kv   *redis.Client
	db   *gorm.DB
	rest *resty.Client
	log  *zap.SugaredLogger
}

// New create a job service instance
func New(kv *redis.Client, db *gorm.DB, rest *resty.Client, log *zap.SugaredLogger) *Service {
	var s = &Service{
		ctx:  context.Background(),
		kv:   kv,
		db:   db,
		rest: rest,
		log:  log,
	}
	return s
}
