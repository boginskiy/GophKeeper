package config

import "github.com/boginskiy/GophKeeper/client/internal/logg"

// Data default.
const (
	APPNAME = "gophclient"
	DESC    = "HI man! It is special CLI application for your computer"
	VERS    = "1.1.01"
	CONFIG  = "config.json"
)

type Config interface {
	GetPortServerGRPC() string
}

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
