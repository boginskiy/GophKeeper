package user

import "github.com/boginskiy/GophKeeper/client/internal/model"

type User interface {
	TakeSystemInfoCurrentUser() (username, uid string)
	ReceiveMess() (string, error)
	SaveLocalUser(*model.User)
}
