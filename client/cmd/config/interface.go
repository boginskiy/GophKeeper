package config

type Config interface {
	GetWaitingTimeResponse() int
	GetCountRetryRequest() int
	GetPortServerGRPC() string
	GetAttempts() int

	GetNameApp() string
	GetDescApp() string
	GetVersionApp() string
	GetConfigFile() string
}
