package config

import (
	"time"

	"github.com/boginskiy/GophKeeper/client/internal/logg"
)

type Conf struct {
	Logg           logg.Logger
	PortServerGRPC string
	Attempts       int
	WaitTimeRes    int
	RetryReq       int

	// Info App
	NameApp    string
	DescApp    string
	VersionApp string
	ConfigFile string
}

func NewConf(
	logger logg.Logger,
	port string,
	attempts int,
	waitTimeRes int,
	retryReq int,

	// Info App
	nameApp string,
	descApp string,
	versionApp string,
	configFile string,

) *Conf {

	return &Conf{
		Logg:           logger,
		PortServerGRPC: port,
		Attempts:       attempts,
		WaitTimeRes:    waitTimeRes,
		RetryReq:       retryReq,

		// Info App
		NameApp:    nameApp,
		DescApp:    descApp,
		VersionApp: versionApp,
		ConfigFile: configFile,
	}
}

func (c *Conf) GetPortServerGRPC() string {
	return c.PortServerGRPC
}

func (c *Conf) GetAttempts() int {
	return c.Attempts
}

func (c *Conf) GetWaitingTimeResponse() int {
	return c.WaitTimeRes * int(time.Millisecond)
}

func (c *Conf) GetCountRetryRequest() int {
	return c.RetryReq
}

func (c *Conf) GetNameApp() string {
	return c.NameApp
}

func (c *Conf) GetDescApp() string {
	return c.DescApp
}

func (c *Conf) GetVersionApp() string {
	return c.VersionApp
}

func (c *Conf) GetConfigFile() string {
	return c.ConfigFile
}
