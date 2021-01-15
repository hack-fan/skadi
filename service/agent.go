package service

import (
	"fmt"

	"github.com/rs/xid"

	"github.com/hack-fan/skadi/types"
)

func (s *Service) AgentAdd(uid string) (*types.Agent, error) {
	var agent = &types.Agent{
		ID:     xid.New().String(),
		UserID: uid,
		Secret: xid.New().String(),
	}
	err := s.db.Create(agent).Error
	if err != nil {
		return nil, fmt.Errorf("create new agent to db failed: %w", err)
	}
	return agent, nil
}
