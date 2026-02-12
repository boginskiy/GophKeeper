package auth

import (
	"context"

	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/internal/errs"
	"github.com/boginskiy/GophKeeper/server/internal/infra"
	"github.com/boginskiy/GophKeeper/server/internal/logg"
	"github.com/boginskiy/GophKeeper/server/internal/model"
	"github.com/boginskiy/GophKeeper/server/internal/repo"
	"github.com/boginskiy/GophKeeper/server/internal/rpc"
	"github.com/boginskiy/GophKeeper/server/pkg"
)

type Auth struct {
	Cfg        config.Config
	Logger     logg.Logger
	JWTService JWTer
	Repo       repo.RepoCreateReadUpdater[*model.User]
}

func NewAuth(
	config config.Config,
	logger logg.Logger,
	jwtService JWTer,
	repo repo.RepoCreateReadUpdater[*model.User]) *Auth {

	return &Auth{Cfg: config, Logger: logger, JWTService: jwtService, Repo: repo}
}

func (a *Auth) Recovery(ctx context.Context, mod *model.User) (token string, err error) {
	// Check user in DB.
	record, err := a.Repo.UpdateRecord(context.Background(), mod)
	if err != nil {
		return "", err
	}

	// Create new token
	infoToken := NewExtraInfoToken(record.ID, record.Email, record.PhoneNumber)
	token, err = a.JWTService.CreateJWT(infoToken)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *Auth) Authentication(ctx context.Context, mod *model.User) (token string, err error) {
	// Check user in DB.
	record, err := a.Repo.ReadRecord(context.Background(), &model.User{Email: mod.Email})
	if err != nil {
		return "", err
	}

	// Check password
	if !pkg.CompareHashAndPassword(record.Password, mod.Password) {
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

func (a *Auth) Registration(ctx context.Context, mod *model.User) (token string, err error) {
	// Create new record with user.
	record, err := a.Repo.CreateRecord(context.Background(), mod)
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
	_, ok := req.(*rpc.RegistRequest)
	return ok
}

// CheckPathToAuth check, if client go to Authentication.
func (a *Auth) CheckPathToAuth(ctx context.Context, req any) bool {
	_, ok := req.(*rpc.AuthRequest)
	return ok
}

// CheckPathToRec check, if client go to Recovery.
func (a *Auth) CheckPathToRec(ctx context.Context, req any) bool {
	_, ok := req.(*rpc.RecovRequest)
	return ok
}
