package user

import (
	"io"
	"os"

	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/utils"
)

type Identity struct {
	Logger       logg.Logger
	FileHendler  utils.FileHandler
	PathToConfig string
}

func NewIdentity(logger logg.Logger, fileHndl utils.FileHandler) *Identity {
	// Path to config file.
	path, err := fileHndl.CreatePathToConfig(config.APPNAME, config.CONFIG)
	logger.CheckWithFatal(err, "error in creating path to config file")

	// Create folder for config file.
	err = fileHndl.CreateFolder(path, 0755)
	logger.CheckWithFatal(err, "error in creating path to config file")

	tmp := &Identity{
		Logger:       logger,
		FileHendler:  fileHndl,
		PathToConfig: path,
	}

	return tmp
}

func (i *Identity) ReadConfigFile(path string, mod os.FileMode) ([]byte, error) {
	file, err := i.FileHendler.ReadOrCreateFile(path, mod)
	if err != nil {
		return []byte{}, err
	}

	defer file.Close()
	return io.ReadAll(file)
}

func (i *Identity) TakePreviosUser() (*model.User, error) {
	// Read config.file.
	dataByte, err := i.ReadConfigFile(i.PathToConfig, 0755)
	if err != nil {
		i.Logger.RaiseError(err, "error on take info from config file", nil)
		return nil, err
	}

	// If config file is empty.
	if len(dataByte) == 0 {
		return nil, ErrEmptyConfigFile
	}

	// Deserialization.
	previosUser := &model.User{}
	err = i.FileHendler.Deserialization(dataByte, previosUser)
	if err != nil {
		i.Logger.RaiseError(err, "error deserialization info from config file", nil)
		return nil, err
	}
	return previosUser, nil
}

func (i *Identity) SystemIdentification(userCLI *UserCLI) bool {
	previosUser, err := i.TakePreviosUser()
	if err != nil {
		return false
	}

	systemUserName, systemUserId := userCLI.TakeSystemInfoCurrentUser()

	if previosUser.SystemUserName != systemUserName || previosUser.SystemUserId != systemUserId {
		return false
	}

	// Save restored data in userCLI.
	userCLI.User = previosUser
	return true
}

// PutCurrentUser need for save data about current user.
func (i *Identity) SaveCurrentUser(user *model.User) error {
	file, err := i.FileHendler.TruncateFile(i.PathToConfig, 0755)
	if err != nil {
		i.Logger.RaiseError(err, "error in trancate config file", nil)
		return err
	}

	dataByte, err := i.FileHendler.Serialization(user)
	if err != nil {
		i.Logger.RaiseError(err, "error in serialization config file", nil)
		return err
	}

	_, err = file.Write(dataByte)
	if err != nil {
		i.Logger.RaiseError(err, "error in write config file", nil)
		return err
	}

	defer file.Close()
	return nil
}
