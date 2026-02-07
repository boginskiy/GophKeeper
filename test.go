package main

import "github.com/boginskiy/GophKeeper/client/pkg"

const (
	GmailSMTPHost = "smtp.gmail.com"
	GmailSMTPPort = "587"
	EmailFrom     = "gophkeeper@gmail.com"
	AppPassword   = "upiplnvviujgnevc"
)

// smtp.gmail.com
// 587
// gophkeeper@gmail.com
// upiplnvviujgnevc

func main() {
	// cfg := config.NewArgsCLI(&TestLogg{})

	sender := pkg.NewEmailSend(GmailSMTPHost, GmailSMTPPort, EmailFrom, AppPassword)

	sender.SendEmail("1.boginskiy@mail.ru", "Account recovery", "123456")

}
