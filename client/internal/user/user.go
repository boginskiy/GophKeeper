package user

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"time"
)

type UserCLI struct {
	Name    string
	InMess  chan string
	OutMess chan string
	Scanner *bufio.Scanner
}

func NewUserCLI(ctx context.Context, name string, out, in chan string) *UserCLI {
	tmp := &UserCLI{
		Name:    name,
		OutMess: out,
		InMess:  in,

		Scanner: bufio.NewScanner(os.Stdin),
	}

	// tmp.Reciver(ctx)

	return tmp
}

func (u *UserCLI) procInMess(text string) {

}

// func (u *UserCLI) Reciver(ctx context.Context) {
// 	for {
// 		select {
// 		case text := <-u.InMess:
// 			u.procInMess(text)

// 		case <-ctx.Done():
// 			return
// 		}
// 	}
// }

func (u *UserCLI) Reciver() (string, error) {
	if !u.Scanner.Scan() {
		return "", errors.New("error user reciver")
	}
	return u.Scanner.Text(), nil
}

func (u *UserCLI) Sender(text string) {
	u.OutMess <- text
}

func (u *UserCLI) ReceiveMess() (string, error) {
	time.Sleep(500 * time.Millisecond)

	fmt.Fprintf(os.Stdout, "%s: ", u.Name)

	if !u.Scanner.Scan() {
		return "", errors.New("error user reciver")
	}

	return u.Scanner.Text(), nil
}
