package auth

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/pkg"
)

const (
	TOPIC   = "Account recovery"
	LONGNUM = 6
)

type Recovery struct {
	// Ctx          context.Context
	Cfg          config.Config
	Logg         logg.Logger
	MailChan     <-chan string
	CodeChan     chan<- string
	EmailSender  pkg.EmailSender
	RandomNumber string
}

func NewRecovery(
	ctx context.Context,
	cfg config.Config,
	logger logg.Logger,
	mailch <-chan string,
	codeChan chan<- string,
) *Recovery {

	emailSender := pkg.NewEmailSend(
		cfg.GetSMTPHost(),
		cfg.GetSMTPPort(),
		cfg.GetEmailFrom(),
		cfg.GetAppPassword())

	tmp := &Recovery{
		// Ctx:         ctx,
		Cfg:         cfg,
		Logg:        logger,
		MailChan:    mailch,
		EmailSender: emailSender,
		CodeChan:    codeChan,
	}

	go tmp.ReceiveMail(ctx)

	return tmp
}

func (d *Recovery) ReceiveMail(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case mail := <-d.MailChan:
			if d.RecoveryPassword(mail) {
				d.CodeChan <- d.RandomNumber
			}
		}
	}
}

func (d *Recovery) GetRandomNumber() string {
	return d.RandomNumber
}

func (d *Recovery) generatingRandomNumber(long int) string {
	// Создаем источник случайности на основе текущего времени
	source := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(source)
	result := ""

	for i := 0; i < long; i++ {
		digit := strconv.Itoa(rand.Intn(10))
		result += digit
	}
	d.RandomNumber = result
	return result
}

func (d *Recovery) RecoveryPassword(email string) bool {
	// Генерация случайного числа
	d.generatingRandomNumber(LONGNUM)

	// Отправка пользователю на email
	err := d.EmailSender.SendEmail(email, TOPIC, d.RandomNumber)
	if err != nil {
		d.Logg.RaiseError(err, "error about sending email", nil)
		return false
	}
	return true
}
