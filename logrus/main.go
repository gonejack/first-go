package main

import (
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetFormatter(new(logrus.TextFormatter))

	logger := logrus.WithFields(
		logrus.Fields{
			"user_id": 10010,
			"ip":      "192.168.32.15",
		},
	)

	logger.Trace("trace msg")
	logger.Debug("debug msg")
	logger.Info("info msg")
	logger.Warn("warn msg")
	logger.Error("error msg")
	logger.Fatal("fatal msg")
	logger.Panic("panic msg")
}
