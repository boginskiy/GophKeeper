package main

import (
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/app"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
)

const (
	APPNAME = "gophclient"
	DESC    = "HI man! It is special CLI application for your computer"
	VERS    = "1.1.01"
	CONFIG  = "config.json"
)

var (
	port        = ":8080"
	attempts    = 3
	waitTimeRes = 500
	retryReq    = 3
)

func main() {
	logg := logg.NewLogg("main_client.log", "INFO")
	cfg := config.NewConf(logg, port, attempts, waitTimeRes, retryReq, APPNAME, DESC, VERS, CONFIG)
	app.NewApp(cfg, logg).Init()
}
