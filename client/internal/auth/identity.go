package auth

import (
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/errs"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"github.com/boginskiy/GophKeeper/client/internal/utils"
)

type Identity struct {
	Cfg          config.Config
	Logger       logg.Logger
	FileHendler  utils.FileHandler
	PathHandler  utils.PathHandler
	PathToConfig string
}

func NewIdentity(
	cfg config.Config,
	logger logg.Logger,
	fileHandler utils.FileHandler,
	pathHandler utils.PathHandler) *Identity {

	// Path to config file.
	path, err := pathHandler.CreatePathToConfig(config.APPNAME, config.CONFIG)
	logger.CheckWithFatal(err, "error in creating path to config file")

	// Create folder for config file.
	err = fileHandler.CreateFolder(path, 0755)
	logger.CheckWithFatal(err, "error in creating path to config file")

	tmp := &Identity{
		Logger:       logger,
		FileHendler:  fileHandler,
		PathToConfig: path,
	}

	return tmp
}

func (i *Identity) Shutdown(done <-chan struct{}, user user.User) {
	//  Registration interruption.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		defer stop()

		select {
		case <-ctx.Done():
			i.SaveCurrentUser(user)
			return
		case <-done:
			return
		}
	}()
}

func (i *Identity) Identification(user user.User) bool {
	previosUser, err := i.takePreviosUser()
	if err != nil {
		return false
	}
	systemUserName, systemUserId := user.GetSystemInfo()
	if previosUser.SystemUserName != systemUserName || previosUser.SystemUserId != systemUserId {
		return false
	}
	// Save restored data in user.
	user.SavePreviosUser(previosUser)

	return true
}

// PutCurrentUser need for save data about current user.
func (i *Identity) SaveCurrentUser(user user.User) {
	file, err := i.FileHendler.TruncateFile(i.PathToConfig, 0755)
	if err != nil {
		i.Logger.RaiseError(err, "error in trancate config file", nil)
		return
	}

	dataByte, err := utils.Serialization(user.GetModelUser())
	if err != nil {
		i.Logger.RaiseError(err, "error in serialization config file", nil)
		return
	}

	_, err = file.Write(dataByte)
	if err != nil {
		i.Logger.RaiseError(err, "error in write config file", nil)
		return
	}

	defer file.Close()
}

func (i *Identity) readConfigFile(path string, mod os.FileMode) ([]byte, error) {
	file, err := i.FileHendler.ReadOrCreateFile(path, mod)
	if err != nil {
		return []byte{}, err
	}

	defer file.Close()
	return io.ReadAll(file)
}

func (i *Identity) takePreviosUser() (*model.User, error) {
	// Read config.file.
	dataByte, err := i.readConfigFile(i.PathToConfig, 0755)
	if err != nil {
		i.Logger.RaiseError(err, "error on take info from config file", nil)
		return nil, err
	}

	// If config file is empty.
	if len(dataByte) == 0 {
		return nil, errs.ErrEmptyConfigFile
	}

	// Deserialization.
	previosUser := &model.User{}
	err = utils.Deserialization(dataByte, previosUser)
	if err != nil {
		i.Logger.RaiseError(err, "error deserialization info from config file", nil)
		return nil, err
	}
	return previosUser, nil
}
