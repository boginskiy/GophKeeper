package tests

var (
	port        = ":8080"
	attempts    = 3
	waitTimeRes = 500
	retryReq    = 3
)

const (
	APPNAME = "gophclient"
	DESC    = "HI man! It is special CLI application for your computer"
	VERS    = "1.1.01"
	CONFIG  = "test_config.json"
)

// TODO!
// Добавить бы автоматическое подниятие сервера, далее прогон тестов
// далее тушим сервер. Так же подумать о тестовой БД
