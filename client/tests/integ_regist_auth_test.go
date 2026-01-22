package tests

import (
	"fmt"
	"testing"

	"github.com/boginskiy/GophKeeper/client/internal/api"
	"github.com/boginskiy/GophKeeper/client/internal/auth"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func testRegistration(t *testing.T, service api.ServiceAPI) {
	tests := []struct {
		name string
		user model.User
		code codes.Code
	}{
		{"test successful registration",
			model.User{UserName: "Dima", Email: "email", PhoneNumber: "8980", Password: "1234"},
			codes.OK},

		{"test error of creating new user. Empty password",
			model.User{UserName: "Dima_1", Email: "email_1", PhoneNumber: "8980", Password: ""},
			codes.InvalidArgument},

		{"test error not unique email",
			model.User{UserName: "Dima_2", Email: "email", PhoneNumber: "8980", Password: "1234"},
			codes.AlreadyExists},
	}

	for _, tt := range tests {
		_, err := service.Registration(tt.user)

		// Преобразование ошибок в коды
		code := auth.ModifyErrServerOnCode(err)

		// print
		fmt.Println(">>>", tt.code, code)

		assert.Equal(t, tt.code, code)
	}
}

func testAuthentication(t *testing.T, service api.ServiceAPI) {
	tests := []struct {
		name string
		user model.User
		code codes.Code
	}{
		{"test successful authentication",
			model.User{Email: "email", Password: "1234"},
			codes.OK},

		{"test error of authentication user. Bad password",
			model.User{Email: "email", Password: "5678"},
			codes.Unauthenticated},

		{"test error of authentication user. email Not found.",
			model.User{Email: "emailBad", Password: "1234"},
			codes.NotFound},
	}

	for _, tt := range tests {
		_, err := service.Authentication(tt.user)

		// Преобразование ошибок в коды
		code := auth.ModifyErrServerOnCode(err)

		// print
		fmt.Println(">>>", tt.code, code)

		assert.Equal(t, tt.code, code)
	}
}
