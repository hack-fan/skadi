package service

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"

	"github.com/hack-fan/skadi/types"
)

func agentAuthKey(key string) string {
	return "agent:auth:" + key
}

func (s *Service) AuthValidator(key string, c echo.Context) (bool, error) {
	var aid string
	var err error
	// find in redis
	aid, err = s.kv.Get(s.ctx, agentAuthKey(key)).Result()
	if err == redis.Nil {
		// find in db
		agent := new(types.Agent)
		err = s.db.First(agent, "secret = ?", key).Error
		if gorm.IsRecordNotFoundError(err) {
			// invalid key
			return false, nil
		} else if err != nil {
			return false, fmt.Errorf("err read key from db: %w", err)
		}
		aid = agent.ID
		// save in redis
		err = s.kv.Set(s.ctx, agentAuthKey(key), aid, 24*time.Hour).Err()
		if err != nil {
			go s.notify(err)
		}
	} else if err != nil {
		return false, fmt.Errorf("err read key from redis: %w", err)
	}
	// set context
	c.Set("aid", aid)

	return true, nil
}
