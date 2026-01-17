package client

type Requester interface {
	Request(string)
	Run()
}
