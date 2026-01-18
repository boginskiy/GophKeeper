package main

import (
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/app"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
)

func main() {
	logg := logg.NewLogg("main.log", "INFO")
	cfg := config.NewConf(logg)
	app.NewApp(cfg, logg).Run()
}

// TODO...

// Сделать приветствие, сделать команду help и вывод всех доступных команд
// Пользователь проходит процедуру первичной регистрации.
// Пользователь добавляет в клиент новые данные.
// Клиент синхронизирует данные с сервером.

// Пользователь проходит процедуру аутентификации.
// Пользователь запрашивает данные.
// Клиент отображает данные для пользователя.

// TODO...
// На клиенте должны быть данные по пользователю. Сохранить в файл как вариант. Предусмотреть обновление данных
// С эими данными далее будет осуществляться идентификация и т.п.
