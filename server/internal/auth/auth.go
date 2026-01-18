package auth

import (
	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/internal/logg"
)

type Auth struct {
	Cfg  config.Config
	Logg logg.Logger
}

func NewAuth(config config.Config, logger logg.Logger) *Auth {
	return &Auth{Cfg: config, Logg: logger}
}
