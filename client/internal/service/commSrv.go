package service

import (
	"context"

	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
)

type CommandService struct {
	Cfg      config.Config
	Logg     logg.Logger
	messChan chan string
}

func NewCommandService(ctx context.Context, cfg config.Config, logger logg.Logger, ch chan string) *CommandService {
	tmp := &CommandService{
		Cfg:      cfg,
		Logg:     logger,
		messChan: ch,
	}

	go tmp.ReceiveMess(ctx)

	return tmp
}

func (c *CommandService) ReceiveMess(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case mess := <-c.messChan:
			c.Router(mess)
		}
	}
}

func (c *CommandService) Router(mess string) {

}

// TODO...
// Что я хочу хранить на удаленном сервере ?
// Как я это туда буду передавать ?
// Как я это буду от туда забирать ?
