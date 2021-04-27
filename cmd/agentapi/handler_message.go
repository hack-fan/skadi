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
	_, err = h.evm.Pub(&types.Message{
		ID:        xid.New().String(),
		AgentID:   aid,
		UserID:    uid,
		Type:      types.MessageTypeInfo,
		Message:   input.Message,
		CreatedAt: time.Now(),
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
	_, err = h.evm.Pub(&types.Message{
		ID:        xid.New().String(),
		AgentID:   aid,
		UserID:    uid,
		Type:      types.MessageTypeWarning,
		Message:   input.Message,
		CreatedAt: time.Now(),
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
	_, err = h.evm.Pub(&types.Message{
		ID:        xid.New().String(),
		AgentID:   aid,
		UserID:    uid,
		Type:      types.MessageTypeText,
		Message:   input.Message,
		CreatedAt: time.Now(),
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
	_, err = h.evm.Pub(&types.Message{
		ID:        xid.New().String(),
		AgentID:   aid,
		UserID:    uid,
		Type:      types.MessageTypeAuto,
		Message:   input.Message,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}
	return c.NoContent(202)
}
