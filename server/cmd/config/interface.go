package config

type Config interface {
	GetConnDB() string
	GetServerGrpc() string
	GetTokenSecretKey() string
	GetTokenLifetime() int
}
