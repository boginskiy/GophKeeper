package tests

import (
	"fmt"
	"testing"

	"github.com/boginskiy/GophKeeper/client/internal/auth"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func testAuthService(t *testing.T, srv auth.Auth, user user.User) {
	testRegistration(t, srv, user)
	testAuthentication(t, srv, user)
}

func testRegistration(t *testing.T, service auth.Auth, user user.User) {
	tests := []struct {
		name    string
		modUser *model.User
		code    codes.Code
	}{
		{"test successful registration",
			&model.User{UserName: "Tester2", Email: "tester2@mail.ru", PhoneNumber: "89109109911", Password: "1234"},
			codes.OK},

		{"test error of creating new user. Empty password",
			&model.User{UserName: "Tester_1", Email: "tester@mail.ru_1", PhoneNumber: "89109109910", Password: ""},
			codes.InvalidArgument},

		{"test error not unique email",
			&model.User{UserName: "Tester_2", Email: "tester@mail.ru", PhoneNumber: "89109109910", Password: "1234"},
			codes.AlreadyExists},
	}

	for i, tt := range tests {
		_, err := service.Registration(user, tt.modUser)

		// Преобразование ошибок в коды
		code := auth.ModifyErrServerOnCode(err)

		// print
		fmt.Println(">>>", i, tt.code, code)
		assert.Equal(t, tt.code, code)
	}
}

func testAuthentication(t *testing.T, service auth.Auth, user user.User) {
	tests := []struct {
		name string
		user *model.User
		code codes.Code
	}{
		{"test successful authentication",
			&model.User{Email: "tester@mail.ru", Password: "1234"},
			codes.OK},

		{"test error of authentication user. Bad password",
			&model.User{Email: "tester@mail.ru", Password: "5678"},
			codes.Unauthenticated},

		{"test error of authentication user. email Not found.",
			&model.User{Email: "emailBad", Password: "1234"},
			codes.NotFound},
	}

	for _, tt := range tests {
		_, err := service.Authentication(user, tt.user)

		// Преобразование ошибок в коды
		code := auth.ModifyErrServerOnCode(err)

		// print
		fmt.Println(">>>", tt.code, code)

		assert.Equal(t, tt.code, code)
	}
}
