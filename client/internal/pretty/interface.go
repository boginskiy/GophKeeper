package pretty

type Looker interface {
	PrintInfo(whoSay, info string)
	Hello(requestor, responder string)
}
