package utils

import (
	"github.com/sirupsen/logrus"
	"os"
)

func NewLog(path string) *logrus.Logger {
	log := logrus.New()

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logrus.Error(err)
	} else {
		log.Out = file
	}
	log.SetLevel(logrus.InfoLevel)

	log.SetFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	return log
}
