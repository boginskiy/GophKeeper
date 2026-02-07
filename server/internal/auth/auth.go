package auth

import (
	"context"

	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/internal/errs"
	"github.com/boginskiy/GophKeeper/server/internal/infra"
	"github.com/boginskiy/GophKeeper/server/internal/logg"
	"github.com/boginskiy/GophKeeper/server/internal/model"
	repo "github.com/boginskiy/GophKeeper/server/internal/repo"
	"github.com/boginskiy/GophKeeper/server/internal/rpc"
	"github.com/boginskiy/GophKeeper/server/pkg"
)

type Auth struct {
	Cfg        config.Config
	Logg       logg.Logger
	JWTService JWTer
	Repo       repo.RepoCreateReadUpdater[*model.User]
}

func NewAuth(
	config config.Config,
	logger logg.Logger,
	jwtSrv JWTer,
	repo repo.RepoCreateReadUpdater[*model.User]) *Auth {

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

func (a *Auth) Recovery(ctx context.Context, req *rpc.RecoverUserRequest) (token string, err error) {
	updateUser, err := NewUser("", req.Email, req.Password, "")
	if err != nil {
		return "", errs.ErrCreateUser.Wrap(err)
	}

	// Check user in DB.
	record, err := a.Repo.UpdateRecord(context.Background(), updateUser)
	if err != nil {
		return "", errs.ErrUpdateUser.Wrap(err)
	}

	// Create new token
	infoToken := NewExtraInfoToken(record.ID, record.Email, record.PhoneNumber)
	token, err = a.JWTService.CreateJWT(infoToken)
	if err != nil {
		return "", errs.ErrCreateToken.Wrap(err)
	}

	return token, nil
}

func (a *Auth) Authentication(ctx context.Context, req *rpc.AuthUserRequest) (token string, err error) {
	// Check user in DB.
	record, err := a.Repo.ReadRecord(context.Background(), &model.User{Email: req.Email})
	if err != nil {
		return "", err
	}

	// Check password
	if !pkg.CompareHashAndPassword(record.Password, req.Password) {
		return "", errs.ErrUserPassword
	}

	// Create new token
	infoToken := NewExtraInfoToken(record.ID, record.Email, record.PhoneNumber)
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
	record, err := a.Repo.CreateRecord(context.Background(), newUser)
	if err != nil {
		return "", errs.ErrCreateUser.Wrap(err)
	}

	// Create token
	infoToken := NewExtraInfoToken(record.ID, record.Email, record.PhoneNumber)

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

	// Check, if client go to Recovery.
	if a.CheckPathToRec(ctx, req) {
		return nil, true
	}

	token := infra.TakeDataFromCtx(ctx, "authorization")

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

// CheckPathToRec check, if client go to Recovery.
func (a *Auth) CheckPathToRec(ctx context.Context, req any) bool {
	_, ok := req.(*rpc.RecoverUserRequest)
	return ok
}
