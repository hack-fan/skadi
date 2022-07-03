package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hack-fan/x/xerr"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/hack-fan/skadi/types"
)

func agentAuthKey(key string) string {
	return "agent:auth:" + key
}

func (s *Service) clearAgentAuthCache(key string) {
	err := s.kv.Del(s.ctx, agentAuthKey(key)).Err()
	if err != nil {
		s.log.Errorf("clear agent secret in redis failed:%s", err)
		return
	}
}

func (s *Service) AuthValidator(key string, c echo.Context) (bool, error) {
	var aid, uid string
	// find in redis
	val, err := s.kv.Get(s.ctx, agentAuthKey(key)).Result()
	if err == redis.Nil {
		// find in db
		agent := new(types.Agent)
		err = s.db.First(agent, "secret = ?", key).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// invalid key
			return false, nil
		} else if err != nil {
			return false, fmt.Errorf("err read key from db: %w", err)
		}
		// unavailable
		if !agent.Available {
			return false, xerr.New(400, "UnavailableAgent", "agent not available")
		}
		// ok
		aid = agent.ID
		uid = agent.UserID
		// save in redis
		err = s.kv.Set(s.ctx, agentAuthKey(key), aid+","+uid, 24*time.Hour).Err()
		if err != nil {
			go s.log.Error(err)
		}
	} else if err != nil {
		return false, fmt.Errorf("err read key from redis: %w", err)
	} else {
		// parse val
		pair := strings.Split(val, ",")
		if len(pair) != 2 {
			return false, fmt.Errorf("parse cached agent auth failed: %w", err)
		}
		aid = pair[0]
		uid = pair[1]
	}
	// set context
	c.Set("aid", aid)
	c.Set("uid", uid)

	return true, nil
}
