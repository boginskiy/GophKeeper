package client

import (
	"context"
	"fmt"
	"os"
	"time"
)

const NAME = "CLIENT"

type ClientCLI struct {
	Name    string
	InMess  chan string
	OutMess chan string
}

func NewClientCLI(ctx context.Context, in, out chan string) *ClientCLI {
	return &ClientCLI{
		Name:    NAME,
		InMess:  in,
		OutMess: out,
	}
}

func (c *ClientCLI) Reciver() (string, error) {
	return "", nil
}

func (c *ClientCLI) Sender(text string) {
	fmt.Fprintf(os.Stdout, "%s: %s\n\r", c.Name, text)
	// c.OutMess <- text
}

func (c *ClientCLI) Pprint(text string) {
	for _, ch := range text {
		fmt.Fprintf(os.Stdout, "%c", ch)
		time.Sleep(30 * time.Millisecond)
	}
}

func (c *ClientCLI) SendMess(text ...string) {
	for _, t := range text {
		c.Pprint(fmt.Sprintf("%s: %s\n\r", c.Name, t))
	}
}
