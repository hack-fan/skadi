package main

import (
	"github.com/hack-fan/x/rdb"
	"github.com/hack-fan/x/xdb"
	"github.com/hack-fan/x/xlog"

	"github.com/hack-fan/skadi/event"
)

// Settings will load from env and docker secret
type Settings struct {
	Debug      bool   `default:"false"`
	ListenAddr string `default:":1323"` // only change this in local debug
	Hostname   string

	DB xdb.Config

	Redis rdb.Config

	Event event.Config

	Wework xlog.WeworkSender
}
