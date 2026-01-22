package config

import (
	"time"

	"github.com/boginskiy/GophKeeper/server/internal/logg"
)

type Config interface {
	GetPortServerGRPC() string
	GetTokenLiveTime() int
	GetSecretKey() string
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

func (c *Conf) GetTokenLiveTime() int {
	return 3600 * int(time.Second)
}

func (c *Conf) GetSecretKey() string {
	return "Ld5pS4Gw"
}
