package comm

import (
	"fmt"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/internal/cli"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type Root struct {
	Dialoger  cli.Dialoger
	CommText  Commander
	CommBytes Commander
	CommImage Commander
	CommSound Commander
}

func NewRoot(
	dialoger cli.Dialoger,
	commtext Commander,
	commbytes Commander,
	commimage Commander,
	commsound Commander,
) *Root {

	return &Root{
		Dialoger:  dialoger,
		CommText:  commtext,
		CommBytes: commbytes,
		CommImage: commimage,
		CommSound: commsound,
	}
}

func (r *Root) Execute(authOK bool, client *client.ClientCLI, user *user.UserCLI) {
authLoop:
	for authOK {
		// Get command.
		comm, _ := r.Dialoger.GetSomeThing(client, user,
			fmt.Sprintf("%s\n\r%s",
				"Enter the data type you want to work with: \n\r\t text \n\r\t bytes \n\r\t image \n\r\t sound",
				"go out: exit, need help: help"))

		switch comm {
		case "exit", "help":
			break authLoop

		case "text":
			r.CommText.Registration(client, user)
		case "bytes":
			r.CommBytes.Registration(client, user)
		case "image":
			r.CommImage.Registration(client, user)
		case "sound":
			r.CommSound.Registration(client, user)

		default:
			r.Dialoger.ShowIt(client, "Invalid command. Try again...")
		}
	}
}

// пары логин/пароль;
// произвольные текстовые данные;
// произвольные бинарные данные;
// данные банковских карт.

// как забирать данные, как их вносить ? Может отправка в поток ввода ?

// TODO...
// var CommStore = make(map[string]func(), 10)
// Что я хочу хранить на удаленном сервере ?
// Как я это туда буду передавать ?
// Как я это буду от туда забирать ?

// Для хранения музыки и фото буду использовать обычную файловую систему
// Под чтение и запись чего то в файловую систему создаем отдельный поток!
// Применяйте очереди заданий (RabbitMQ, Kafka) для снижения пиковых нагрузок.

// безопасно хранить логины, пароли, бинарные данные и прочую приватную информацию
// Передача пароля без явного вида. Как то надо хешировать ?
// обеспечение безопасности передачи и хранения данных

// типы хранимой инфы
// пары логин/пароль;
// произвольные текстовые данные;
// произвольные бинарные данные;
// данные банковских карт.

// Для любых данных должна быть возможность хранения произвольной текстовой
// метаинформации (принадлежность данных к веб-сайту, личности или банку,
// 	списки одноразовых кодов активации и прочее).

// Пользователь добавляет в клиент новые данные.
// Клиент синхронизирует данные с сервером.
