package service

import "github.com/hack-fan/skadi/types"

func (s *Service) UserSetting(uid string) (*types.UserSetting, error) {
	var us = new(types.UserSetting)
	err := s.db.First(us, "user_id = ?", uid).Error
	if err != nil {
		return nil, err
	}
	return us, nil
}
