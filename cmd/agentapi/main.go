package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"github.com/hack-fan/config"
	"github.com/hack-fan/x/xdb"
	"github.com/hack-fan/x/xecho"
	"github.com/hack-fan/x/xerr"
	"github.com/hack-fan/x/xlog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/hack-fan/skadi/event"
	"github.com/hack-fan/skadi/service"
	"github.com/hack-fan/skadi/types"
)

func main() {
	var err error
	// load config
	var settings = new(Settings)
	config.MustLoad(settings)

	// logger
	var logger = xlog.New(settings.Debug, settings.Wework)
	defer logger.Sync() // nolint
	var log = logger.Sugar()

	// kv
	var kv = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", settings.Redis.Host, settings.Redis.Port),
		Password: settings.Redis.Password,
		DB:       settings.Redis.DB,
	})

	// db
	xdb.SetLogger(log)
	var db = xdb.New(settings.DB)
	if settings.Debug {
		db = db.Debug()
	}
	// auto create table
	go db.AutoMigrate(&types.Job{}, &types.Agent{}) // nolint

	// http client
	var rest = resty.New().SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(60 * time.Second)

	// service
	var s = service.New(kv, db, rest, log)

	// default event center is just log the events
	// if you have event worker, set it to redis in settings
	var ev = event.NewEventCenter(log, settings.Event)

	// handler
	var h = NewHandler(s, ev)

	// Echo instance
	e := echo.New()
	if settings.Debug {
		e.Debug = true
	}
	// Error handler
	e.HTTPErrorHandler = xerr.ErrorHandler
	// Middleware
	e.Use(xecho.LoggerMid())
	e.Use(middleware.Recover())

	// Auth group
	var a = e.Group("", middleware.KeyAuth(s.AuthValidator))

	// Routes
	e.GET("/status", getStatus)

	a.GET("/agent/job", h.GetJob)
	a.PUT("/agent/jobs/:id/succeed", h.PutJobSucceed)
	a.PUT("/agent/jobs/:id/fail", h.PutJobFail)
	a.PUT("/agent/jobs/:id/running", h.PutJobRunning)
	a.POST("/agent/info", h.PostInfo)
	a.POST("/agent/warning", h.PostWarning)
	a.POST("/agent/kf", h.PostText)

	// Start server
	go func() {
		log.Warnf("Agent API Start: %s", settings.Hostname)
		err := e.Start(settings.ListenAddr)
		if err != nil && err != http.ErrServerClosed {
			log.Errorf("Agent API Force Shutting down: %s %s", settings.Hostname, err)
			log.Fatal("force shutting down the server")
		}
	}()
	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = e.Shutdown(ctx)
	if err != nil {
		log.Errorf("Agent API graceful shutdown failed: %s %s", settings.Hostname, err)
		log.Fatal(err)
	}
	log.Warnf("Agent API graceful shutdown: %s", settings.Hostname)
}
