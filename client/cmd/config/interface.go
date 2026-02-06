package config

type SenderEmailConfig interface {
	GetEmailFrom() string
	GetAppPassword() string
	GetSMTPHost() string
	GetSMTPPort() string
}

type Config interface {
	GetReqRetries() int
	GetServerGrpc() string
	GetResTimeout() int
	GetMaxRetries() int

	SenderEmailConfig
}
