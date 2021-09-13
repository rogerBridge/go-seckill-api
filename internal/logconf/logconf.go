package logconf

import (
	"io"
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

// func init() {
// 	// Log as JSON instead of the default ASCII formatter.
// 	logrus.SetFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05.999"})

// 	// Output to stdout instead of the default stderr
// 	// Can be any io.Writer, see below for File example
// 	logrus.SetOutput(os.Stdout)

// 	// Only log the warning severity or above.
// 	logrus.SetLevel(logrus.InfoLevel)
// }

// var BaseLogger = logrus.WithFields(logrus.Fields{
// 	"app":    "go-seckill",
// 	"author": "Roger Bridge",
// })

// // 需要在设定的logger里面添加新的fields时使用
// // 例如: LogWithMethodName(logger, "hello") ==> logger.WithFields(logrus.Fields{"methodName": "hello"})
// func LogWithMethodName(baseLogger *logrus.Entry, methodName string) *logrus.Entry {
// 	r := baseLogger.WithFields(logrus.Fields{
// 		"methodName": methodName,
// 	})
// 	return r
// }

func init() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.WarnLevel)
	mw := io.MultiWriter(os.Stdout, file)
	logrus.SetOutput(mw)
}

var BaseLogger = logrus.WithFields(logrus.Fields{
	"app": "go-seckill",
})
