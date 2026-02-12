package user

import "github.com/boginskiy/GophKeeper/client/internal/model"

type Saver interface {
	SavePreviosUser(previosUser *model.User)
	SaveLocalUser(localUser *model.User)
}

type Getter interface {
	GetSystemInfo() (sysName, sysUid string)
	GetModelUser() *model.User
}

type User interface {
	Getter
	Saver
	ReceiveMess() (string, error)
}
