package client

import (
	"context"
	"fmt"
	"os"
	"time"
)

type ClientCLI struct {
	Name    string
	InMess  chan string
	OutMess chan string
}

func NewClientCLI(ctx context.Context, name string, in, out chan string) *ClientCLI {
	tmp := &ClientCLI{
		Name:    name,
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

// Нужно построить на функции Hello форматированный диалог

func (c *ClientCLI) SendMess(text string) {
	time.Sleep(500 * time.Millisecond)
	fmt.Fprintf(os.Stdout, "%s: %s\n\r", c.Name, text)
}
