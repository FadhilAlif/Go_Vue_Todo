package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()

	// Output ke stdout (terminal)
	Log.SetOutput(os.Stdout)

	// Format: bisa TextFormatter atau JSONFormatter
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Level log default (bisa diganti ke DebugLevel, WarnLevel, dll)
	Log.SetLevel(logrus.InfoLevel)
}
