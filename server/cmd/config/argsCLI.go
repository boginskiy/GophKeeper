package config

import (
	"flag"
	"time"

	"github.com/boginskiy/GophKeeper/server/internal/logg"
)

const (
	DESC    = "HI man! It is special gRPC server for your data"
	APPNAME = "gophserver"
	VERS    = "1.1.01"

	BuildVersion = "N/A"
	BuildDate    = "N/A"
	BuildCommit  = "N/A"
)

type ArgsCLI struct {
	Logg logg.Logger

	ServerGrpc     string // Port for the gRPC server.
	TokenSecretKey string // Token secret key.
	TokenLifetime  int    // Token lifetime.
}

func NewArgsCLI(logger logg.Logger) *ArgsCLI {
	args := &ArgsCLI{Logg: logger}
	args.ParseFlags()

	return args
}

func (a *ArgsCLI) ParseFlags() {
	// -p localhost:8080
	// -s Ld5pS4Gw
	// -t 3600

	flag.StringVar(&a.ServerGrpc, "p", "localhost:8080", "Port for the gRPC server")
	flag.StringVar(&a.TokenSecretKey, "s", "Ld5pS4Gw", "Token lifetime")
	flag.IntVar(&a.TokenLifetime, "t", 3600, "Token lifetime")

	flag.Parse()
}

func (a *ArgsCLI) GetServerGrpc() string {
	return a.ServerGrpc
}

func (a *ArgsCLI) GetTokenSecretKey() string {
	return a.TokenSecretKey
}

func (a *ArgsCLI) GetTokenLifetime() int {
	return a.TokenLifetime * int(time.Second)
}
