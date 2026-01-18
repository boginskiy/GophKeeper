package config

import "github.com/boginskiy/GophKeeper/server/internal/logg"

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
