package event

import (
	"go.uber.org/zap"

	"github.com/hack-fan/skadi/types"
)

type DefaultEventCenter struct {
	log *zap.SugaredLogger
}

func NewDefaultEventCenter(log *zap.SugaredLogger) *DefaultEventCenter {
	return &DefaultEventCenter{log}
}

func (ec *DefaultEventCenter) Pub(e *types.Event) error {
	ec.log.Infof("ignore event: %+v", e)
	return nil
}

func (ec *DefaultEventCenter) Get() (*types.Event, error) {
	return nil, nil
}
