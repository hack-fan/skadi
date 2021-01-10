package job

import (
	"context"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/hack-fan/skadi/types"
)

// Service has all JOB functions
type Service interface {
	// Pop a job to agent
	Pop(aid string) (*types.JobBasic, error)
	// Push a job by user
	Push(job *types.JobInput) error
}

type service struct {
	log *zap.SugaredLogger
	kv  *redis.Client
	db  *gorm.DB
	ctx context.Context
}

// NewService create a job service instance
func NewService() Service {
	var s = &service{
		ctx: context.Background(),
	}
	return s
}
