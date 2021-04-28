package main

import (
	"context"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hack-fan/config"
	"github.com/hack-fan/x/rdb"
	"github.com/hack-fan/x/xdb"
	"github.com/hack-fan/x/xlog"

	"github.com/hack-fan/skadi/service"
	"github.com/hack-fan/skadi/types"
)

func main() {
	// load config
	var settings = new(Settings)
	config.MustLoad(settings)

	// logger
	var logger = xlog.New(settings.Debug, settings.Wework)
	defer logger.Sync() // nolint
	var log = logger.Sugar()

	// kv
	rdb.SetLogger(log)
	var kv = rdb.New(settings.Redis)

	// db
	xdb.SetLogger(log)
	var db = xdb.New(settings.DB)
	if settings.Debug {
		db = db.Debug()
	}
	// auto create table
	go db.AutoMigrate(&types.Job{}) // nolint

	// http client
	var rest = resty.New().SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(60 * time.Second)

	// service
	var s = service.New(kv, db, rest, log)

	// watch redis expire events
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	pubsub := kv.Subscribe(ctx, "__keyevent@0__:expired")
	log.Info("start watching redis key expired event...")

	log.Warnf("Skadi Watcher started: %s", settings.Hostname)
	ticker := time.NewTicker(time.Minute)
LOOP:
	for {
		select {
		case msg := <-pubsub.Channel():
			key := msg.Payload
			log.Debugw("redis key expired", "key", key)
			if strings.HasPrefix(key, "job:wait:") {
				jid := strings.TrimPrefix(key, "job:wait:")
				s.JobExpire(jid)
			} else if strings.HasPrefix(key, "agent:online:") {
				aid := strings.TrimPrefix(key, "agent:online:")
				s.AgentOffline(aid)
			}
		case <-ticker.C:
			s.DelayedJobCheck()
		case <-ctx.Done():
			stop()
			break LOOP
		}
	}
	// wait for async functions, 3s is enough.
	time.Sleep(time.Second * 3)
	log.Warnf("Skadi Watcher Stopped: %s", settings.Hostname)
}
