package comm

import (
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type Rooter interface {
	ExecuteComm(bool, user.User) bool
}

type Commander interface {
	Registration(user.User, string)
}
