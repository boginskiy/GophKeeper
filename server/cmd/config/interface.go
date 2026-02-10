package config

type Config interface {
	GetCryptoSignature() []byte
	GetTokenSecretKey() string
	GetServerGrpc() string
	GetTokenLifetime() int
	GetConnDB() string
}
