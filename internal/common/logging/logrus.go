package logging

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	STDOUT = "STDOUT"
)

func Init() {
	SetFormatter(logrus.StandardLogger())
	logrus.SetLevel(logrus.DebugLevel)
	setOutput(logrus.StandardLogger())
}

func SetFormatter(logger *logrus.Logger) {
	logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyMsg:   "message",
		},
	})
	if isLocal, _ := strconv.ParseBool(os.Getenv("LOCAL_ENV")); isLocal {
		//logger.SetFormatter(&prefixed.TextFormatter{
		//	ForceFormatting: false,
		//})
	}
}


func setOutput(logger *logrus.Logger) {
	if logOutput := viper.GetString("loggingOutput"); logOutput == STDOUT {
		return
	} else {
		file, err := os.OpenFile(logOutput, os.O_CREATE|os.O_RDWR, 0755)
		if err != nil {
			panic(err)
		}
		logger.SetOutput(file)
	}
}