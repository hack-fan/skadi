package main

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/hack-fan/skadi/service"
	"github.com/hack-fan/skadi/types"
)

type Handler struct {
	s *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{
		s: s,
	}
}

func (h *Handler) GetJob(c echo.Context) error {
	// TODO: auth
	aid := c.Get("aid").(string)
	resp, err := h.s.JobPop(aid)
	if err != nil {
		return err
	}
	if resp != nil {
		return c.JSON(200, resp)
	}
	// no job found
	return c.NoContent(204)
}

func (h *Handler) PutJobSucceed(c echo.Context) error {
	id := c.Param("id")
	req := types.JobResult{}
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	h.s.JobSucceed(id, req.Result)
	return c.NoContent(204)
}

func (h *Handler) PutJobFail(c echo.Context) error {
	id := c.Param("id")
	req := types.JobResult{}
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	h.s.JobFail(id, req.Result)
	return c.NoContent(204)
}

// API status
func getStatus(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
