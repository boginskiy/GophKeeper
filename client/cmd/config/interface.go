package config

type Config interface {
	GetPortServerGRPC() string
	GetAttempts() int
}
