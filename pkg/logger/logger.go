package logger

import (
	"sync"

	"github.com/sirupsen/logrus"
)

// Logger represent standard logger structure
type Logger struct {
	*logrus.Logger
}

var instance *Logger
var once sync.Once

// Log is
func Log() *Logger {
	once.Do(func() {
		log := logrus.New()
		formatter := &logrus.TextFormatter{
			FullTimestamp: true,
		}
		log.SetFormatter(formatter)

		instance = &Logger{log}
	})
	return instance
}
