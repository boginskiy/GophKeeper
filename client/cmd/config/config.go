package config

import (
	"time"

	"github.com/boginskiy/GophKeeper/client/internal/logg"
)

// Data default.
const (
	APPNAME = "gophclient"
	DESC    = "HI man! It is special CLI application for your computer"
	VERS    = "1.1.01"
	CONFIG  = "config.json"
)

type Conf struct {
	Logg           logg.Logger
	PortServerGRPC string
}

func NewConf(logger logg.Logger) *Conf {
	return &Conf{Logg: logger}
}

func (c *Conf) GetPortServerGRPC() string {
	return ":8080"
}

func (c *Conf) GetAttempts() int {
	return 3
}

func (c *Conf) GetWaitingTimeResponse() int {
	return 250 * int(time.Millisecond)
}

func (c *Conf) GetCountRetryRequest() int {
	return 3
}
