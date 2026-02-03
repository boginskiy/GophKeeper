package config

type Config interface {
	GetReqRetries() int
	GetServerGrpc() string
	GetResTimeout() int
	GetMaxRetries() int
}
