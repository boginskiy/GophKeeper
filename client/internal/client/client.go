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
	tmp := &ClientCLI{
		Name:    NAME,
		InMess:  in,
		OutMess: out,
	}

	// go tmp.Reciver(ctx)

	return tmp
}

func (c *ClientCLI) procInMess(text string) {

}

// func (c *ClientCLI) Reciver(ctx context.Context) {
// 	for {
// 		select {
// 		case text := <-c.InMess:
// 			c.procInMess(text)

// 		case <-ctx.Done():
// 			return
// 		}
// 	}
// }

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
