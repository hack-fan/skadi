package main

import (
	"github.com/hack-fan/x/rdb"
	"github.com/hack-fan/x/xdb"
	"github.com/hack-fan/x/xlog"
)

// Settings will load from env and docker secret
type Settings struct {
	Debug      bool   `default:"false"`
	ListenAddr string `default:":1323"` // only change this in local debug
	Hostname   string

	DB xdb.Config

	Redis rdb.Config

	Wework xlog.WeworkSender
}
