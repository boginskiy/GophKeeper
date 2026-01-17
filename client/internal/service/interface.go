package service

type Dialoger interface {
	Sender(text string)
	Reciver() (string, error)
}
