package logger

import (
	"github.com/sirupsen/logrus"
)

func LogOutput(message string) {
	logrus.Info(message)
}
