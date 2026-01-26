package auth

import (
	"context"

	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/internal/errs"
	"github.com/boginskiy/GophKeeper/server/internal/logg"
	"github.com/boginskiy/GophKeeper/server/internal/model"
	repo "github.com/boginskiy/GophKeeper/server/internal/repository"
	"github.com/boginskiy/GophKeeper/server/internal/rpc"
	"github.com/boginskiy/GophKeeper/server/pkg"
	"google.golang.org/grpc/metadata"
)

type Auth struct {
	Cfg        config.Config
	Logg       logg.Logger
	JWTService JWTer
	Repo       repo.CreateReader[*model.User]
}

func NewAuth(
	config config.Config,
	logger logg.Logger,
	jwtSrv JWTer,
	repo repo.CreateReader[*model.User]) *Auth {
	return &Auth{Cfg: config, Logg: logger, JWTService: jwtSrv, Repo: repo}
}

func NewUser(name, email, password, phone string) (*model.User, error) {
	hash, err := pkg.GenerateHash(password)
	if err != nil {
		return nil, err
	}

	return &model.User{
		UserName:    name,
		Email:       email,
		Password:    hash,
		PhoneNumber: phone,
	}, nil
}

func (a *Auth) Authentication(ctx context.Context, req *rpc.AuthUserRequest) (token string, err error) {
	// Check user in DB.
	user, err := a.Repo.ReadRecord(&model.User{Email: req.Email})
	if err != nil {
		// TODO!
		// Пользователь по введенному email не найден в БД.
		// Тут можно как-то обыграть ситуацию. предложить альтернативные варианты.
		// Восстановить доступ. Например телефон и т.п.
		// Базовое поведение. Отсутствуешь в системе? Иди на регистрацию...
		return "", err
	}

	// Check password
	if !pkg.CompareHashAndPassword(user.Password, req.Password) {
		// TODO!
		// Предусмотреть альтернативы для восстановления учетки.
		return "", errs.ErrUserPassword
	}

	// Create new token
	infoToken := NewExtraInfoToken(user.Email, user.PhoneNumber)
	token, err = a.JWTService.CreateJWT(infoToken)
	if err != nil {
		return "", errs.ErrCreateToken.Wrap(err)
	}

	return token, nil
}

func (a *Auth) Registration(ctx context.Context, req *rpc.RegistUserRequest) (token string, err error) {
	// Create new user.
	newUser, err := NewUser(req.Username, req.Email, req.Password, req.Phonenumber)
	if err != nil {
		return "", errs.ErrCreateUser.Wrap(err)
	}

	// Create new record with user.
	_, err = a.Repo.CreateRecord(newUser)
	if err != nil {
		return "", err
	}

	// Create token
	infoToken := NewExtraInfoToken(newUser.Email, newUser.PhoneNumber)

	token, err = a.JWTService.CreateJWT(infoToken)
	if err != nil {
		return "", errs.ErrCreateToken.Wrap(err)
	}

	return token, nil
}

func (a *Auth) Identification(ctx context.Context, req any) (*ExtraInfoToken, bool) {
	// Check, if client go to Authentication.
	if a.CheckPathToAuth(ctx, req) {
		return nil, true
	}
	// Check, if client go to Registration.
	if a.CheckPathToReg(ctx, req) {
		return nil, true
	}

	token := a.takeDataFromCtx(ctx, "authorization")

	// Try Authentication.
	infoToken, err := a.JWTService.GetInfoAndValidJWT(token)
	if err != nil {
		return nil, false
	}

	return infoToken, true
}

// CheckPathToReg check, if client go to Registration.
func (a *Auth) CheckPathToReg(ctx context.Context, req any) bool {
	_, ok := req.(*rpc.RegistUserRequest)
	return ok
}

// CheckPathToAuth check, if client go to Authentication.
func (a *Auth) CheckPathToAuth(ctx context.Context, req any) bool {
	_, ok := req.(*rpc.AuthUserRequest)
	return ok
}

func (a *Auth) takeDataFromCtx(ctx context.Context, data string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		val := md.Get(data)
		if len(val) > 0 {
			return val[0]
		}
	}
	return ""
}
