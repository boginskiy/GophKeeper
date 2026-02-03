package config

import (
	"flag"
	"time"

	"github.com/boginskiy/GophKeeper/client/internal/logg"
)

const (
	DESC    = "HI man! It is special CLI application for your computer"
	CONFIG  = "config.json"
	APPNAME = "gophclient"
	VERS    = "1.1.01"
)

type ArgsCLI struct {
	Logg logg.Logger

	ServerGrpc string // Port for the gRPC server.
	MaxRetries int    // Maximum number of attempts to enter data through the CLI terminal.
	ResTimeout int    // Waiting time for a response from the remote server.
	ReqRetries int    // RetryReq // Number of attempts to reject requests.
}

func NewArgsCLI(logger logg.Logger) *ArgsCLI {
	args := &ArgsCLI{Logg: logger}
	args.ParseFlags()

	return args
}

func (a *ArgsCLI) ParseFlags() {
	// -p localhost:8080
	// -t 500
	// -m 1
	// -r 3

	flag.IntVar(&a.MaxRetries, "m", 1, "Maximum number of attempts to enter data through the CLI terminal")
	flag.IntVar(&a.ResTimeout, "t", 500, "Waiting time for a response from the remote server")
	flag.StringVar(&a.ServerGrpc, "p", "localhost:8080", "Port for the gRPC server")
	flag.IntVar(&a.ReqRetries, "r", 3, "Number of attempts to reject requests")

	flag.Parse()
}

func (a *ArgsCLI) GetServerGrpc() string {
	return a.ServerGrpc
}

func (a *ArgsCLI) GetMaxRetries() int {
	return a.MaxRetries
}

func (a *ArgsCLI) GetResTimeout() int {
	return a.ResTimeout * int(time.Millisecond)
}

func (a *ArgsCLI) GetReqRetries() int {
	return a.ReqRetries
}
