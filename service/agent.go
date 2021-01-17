package service

import (
	"fmt"
	"time"

	"github.com/rs/xid"

	"github.com/hack-fan/skadi/types"
)

func agentOnlineKey(aid string) string {
	return "agent:online:" + aid
}

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

// call after every agent job pull
func (s *Service) AgentOnline(aid string) {
	err := s.kv.Set(s.ctx, agentOnlineKey(aid), time.Now().Unix(), 3*time.Minute).Err()
	if err != nil {
		go s.notify(fmt.Errorf("set agent %s online failed: %w", aid, err))
	}
}

// call after the watcher found agent status switch to offline
func (s *Service) AgentOffline(aid string) {
	// clear the agent queue, set all job as expired
}
