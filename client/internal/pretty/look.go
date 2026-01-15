package pretty

import (
	"bufio"
	"fmt"
	"os"
)

type Look struct {
}

func NewLook() *Look {
	return &Look{}
}

func (l *Look) Hello(requestor, responder string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Fprintf(os.Stdout,
		"%s: Hello! Press the 'Enter'...\n\r%s: ",
		requestor, responder)

	for {
		if scanner.Scan() {
			break
		}
		fmt.Fprintf(os.Stdout,
			"%s: Press the 'Enter...'\n\r%s: ",
			requestor, responder)
	}
}

func (l *Look) PrintInfo(whoSay, info string) {
	fmt.Fprintf(os.Stdout, "%s: %s\n\r", whoSay, info)
}

func (l *Look) Help() {
	fmt.Fprint(os.Stdout,
		`command XXX           definition XXX
		command YYY           definition YYY`)
}

// TODO...
// Предусмотреть очистку данных принудительную
