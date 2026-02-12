package logg

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// LogService - custom logger above logrus.
type LogService struct {
	Log      *logrus.Logger
	Desc     *os.File
	NameFile string
}

func NewLogService(nameFile, level string) *LogService {
	tmpDesc := createLogFile(nameFile)              // Create file.
	tmpLogrus := setupLogrus(tmpDesc, LEVEL[level]) // Settings Logrus.

	return &LogService{
		NameFile: nameFile,
		Desc:     tmpDesc,
		Log:      tmpLogrus,
	}
}

func (e *LogService) Close() {
	e.Desc.Close()
}

func (e *LogService) RaiseInfo(msg string, dataMap Fields) {
	fmt.Fprintln(os.Stdout, msg)
	e.Log.WithFields(logrus.Fields(dataMap)).Info(msg)
}

func (e *LogService) RaiseWarn(msg string, dataMap Fields) {
	fmt.Fprintln(os.Stdout, msg)
	e.Log.WithFields(logrus.Fields(dataMap)).Warn(msg)
}

func (e *LogService) RaiseError(err error, msg string, dataMap Fields) {
	if err != nil {
		fmt.Fprintf(os.Stdout, "error: %s, mess: %s\n\r", err.Error(), msg)

		if dataMap != nil {
			e.Log.WithFields(logrus.Fields(dataMap)).Error(msg)
		} else {
			e.Log.WithFields(logrus.Fields(Fields{"error": err.Error()})).Error(msg)
		}
	}
}

func (e *LogService) RaiseFatal(err error, msg string, dataMap Fields) {
	if err != nil {
		fmt.Fprintf(os.Stdout, "error: %s, mess: %s\n\r", err.Error(), msg)

		if dataMap != nil {
			e.Log.WithFields(logrus.Fields(dataMap)).Fatal(msg)
		} else {
			e.Log.WithFields(logrus.Fields(Fields{"fatal": err.Error()})).Fatal(msg)
		}
	}
}

func (e *LogService) RaisePanic(err error, msg string, dataMap Fields) {
	if err != nil {
		fmt.Fprintf(os.Stdout, "error: %s, mess: %s\n\r", err.Error(), msg)
		e.Log.WithFields(logrus.Fields(dataMap)).Panic(msg)
	}
}

func (e *LogService) CheckWithFatal(err error, msg string) {
	if err != nil {
		e.RaiseFatal(err, msg, nil)
	}
}

func (e *LogService) CheckWithError(err error, msg string) {
	if err != nil {
		e.RaiseError(err, msg, nil)
	}
}
