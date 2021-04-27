package service

import (
	"context"

	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"github.com/hack-fan/jq"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	ctx      context.Context
	kv       *redis.Client
	db       *gorm.DB
	rest     *resty.Client
	log      *zap.SugaredLogger
	evm      *jq.Queue
	validate *validator.Validate
}

// New create a job service instance
func New(kv *redis.Client, db *gorm.DB, rest *resty.Client, log *zap.SugaredLogger, evm *jq.Queue) *Service {
	var s = &Service{
		ctx:      context.Background(),
		kv:       kv,
		db:       db,
		rest:     rest,
		log:      log,
		evm:      evm,
		validate: validator.New(),
	}
	return s
}
