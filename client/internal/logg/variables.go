package logg

import "github.com/sirupsen/logrus"

// Fields - struct for extra field for logger.
type Fields map[string]any

// LEVEL - level marker for logger.
var LEVEL = map[string]logrus.Level{
	"DEBUG": logrus.DebugLevel,
	"INFO":  logrus.InfoLevel,
	"WARN":  logrus.WarnLevel,
	"ERROR": logrus.ErrorLevel,
	"FATAL": logrus.FatalLevel,
	"PANIC": logrus.PanicLevel,
}
