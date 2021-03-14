package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"github.com/hack-fan/config"
	"github.com/hack-fan/x/xdb"
	"github.com/hack-fan/x/xecho"
	"github.com/hack-fan/x/xerr"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

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
	var logger *zap.Logger
	if settings.Debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		panic(err)
	}
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
	a.POST("/agent/info", h.PostInfo)
	a.POST("/agent/warning", h.PostWarning)

	// Start server
	e.Logger.Fatal(e.Start(settings.ListenAddr))
}
