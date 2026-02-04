package config

type Config interface {
	GetServerGrpc() string
	GetTokenSecretKey() string
	GetTokenLifetime() int
}

// type Conf struct {
// 	Logg       logg.Logger
// 	ServerGrpc string
// }

// func NewConf(logger logg.Logger) *Conf {
// 	return &Conf{Logg: logger}
// }

// func (c *Conf) GetServerGrpc() string {
// 	return ":8080"
// }

// func (c *Conf) GetTokenLiveTime() int {
// 	return 3600 * int(time.Second)
// }

// func (c *Conf) GetSecretKey() string {
// 	return "Ld5pS4Gw"
// }
