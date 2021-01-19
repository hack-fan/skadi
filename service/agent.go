package service

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
	"github.com/vmihailenco/msgpack/v5"

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
	for {
		var job = new(types.JobBasic)
		data, err := s.kv.RPop(s.ctx, agentQueueKey(aid)).Bytes()
		s.log.Debugw("pop", "data", string(data), "err", err)
		if err == redis.Nil {
			break
		} else if err != nil {
			go s.notify(fmt.Errorf("pop job from queue error: %w", err))
			return
		}
		err = msgpack.Unmarshal(data, job)
		if err != nil {
			s.notify(fmt.Errorf("msgpack unmarshal job basic error: %w", err))
			return
		}
		s.JobExpire(job.ID)
	}
}
