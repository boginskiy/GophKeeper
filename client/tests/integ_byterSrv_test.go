package tests

import (
	"fmt"
	"testing"

	"github.com/boginskiy/GophKeeper/client/internal/service"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

var TESTFILE = "/home/ali/dev/GophKeeper/client/tests/store/test.txt"

func testByterService(t *testing.T, service service.BytesServicer, user user.User) {

	any, err := service.Upload(user, TESTFILE)

	fmt.Printf("%+v, %v\n\r", any, err)

}
