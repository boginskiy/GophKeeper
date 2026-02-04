package client

type Client interface {
	SendMess(text ...string)
}
