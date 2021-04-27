package event

import (
	"context"

	"go.uber.org/zap"

	"github.com/hack-fan/skadi/types"
)

type DefaultEventCenter struct {
	log *zap.SugaredLogger
}

func NewDefaultEventCenter(log *zap.SugaredLogger) *DefaultEventCenter {
	return &DefaultEventCenter{log}
}

func (ec *DefaultEventCenter) Pub(e *types.Message) error {
	ec.log.Infof("ignore event: %+v", e)
	return nil
}

func (ec *DefaultEventCenter) Get() (*types.Message, error) {
	return nil, nil
}

// default event center can not have worker
func (ec *DefaultEventCenter) StartWorker(context.Context, types.EventHandler) {}
