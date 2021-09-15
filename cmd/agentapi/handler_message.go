package main

import (
	"time"

	"github.com/hack-fan/skadi/types"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
)

func (h *Handler) PostInfo(c echo.Context) error {
	aid := c.Get("aid").(string)
	uid := c.Get("uid").(string)
	input := new(types.MessageInput)
	err := c.Bind(input)
	if err != nil {
		return err
	}
	input.Type = types.MessageTypeInfo
	input.Level = types.MessageLevelInfo
	_, err = h.evm.Pub(&types.Message{
		ID:           xid.New().String(),
		AgentID:      aid,
		UserID:       uid,
		MessageInput: *input,
		CreatedAt:    time.Now(),
	})
	if err != nil {
		return err
	}
	return c.NoContent(202)
}

func (h *Handler) PostWarning(c echo.Context) error {
	aid := c.Get("aid").(string)
	uid := c.Get("uid").(string)
	input := new(types.MessageInput)
	err := c.Bind(input)
	if err != nil {
		return err
	}
	input.Type = types.MessageTypeWarning
	input.Level = types.MessageLevelWarning
	_, err = h.evm.Pub(&types.Message{
		ID:           xid.New().String(),
		AgentID:      aid,
		UserID:       uid,
		MessageInput: *input,
		CreatedAt:    time.Now(),
	})
	if err != nil {
		return err
	}
	return c.NoContent(202)
}

func (h *Handler) PostText(c echo.Context) error {
	aid := c.Get("aid").(string)
	uid := c.Get("uid").(string)
	input := new(types.MessageInput)
	err := c.Bind(input)
	if err != nil {
		return err
	}
	input.Type = types.MessageTypeText
	input.Level = types.MessageLevelInfo
	_, err = h.evm.Pub(&types.Message{
		ID:           xid.New().String(),
		AgentID:      aid,
		UserID:       uid,
		MessageInput: *input,
		CreatedAt:    time.Now(),
	})
	if err != nil {
		return err
	}
	return c.NoContent(202)
}

func (h *Handler) PostAuto(c echo.Context) error {
	aid := c.Get("aid").(string)
	uid := c.Get("uid").(string)
	input := new(types.MessageInput)
	err := c.Bind(input)
	if err != nil {
		return err
	}
	input.Type = types.MessageTypeAuto
	input.Level = types.MessageLevelInfo
	_, err = h.evm.Pub(&types.Message{
		ID:           xid.New().String(),
		AgentID:      aid,
		UserID:       uid,
		MessageInput: *input,
		CreatedAt:    time.Now(),
	})
	if err != nil {
		return err
	}
	return c.NoContent(202)
}
