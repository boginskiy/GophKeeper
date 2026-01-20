package config

type Config interface {
	GetWaitingTimeResponse() int
	GetCountRetryRequest() int
	GetPortServerGRPC() string
	GetAttempts() int
}
