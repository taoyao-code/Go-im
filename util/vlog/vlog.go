package vlog

import (
	"os"

	"github.com/sirupsen/logrus"
)

var ErrorLog *logrus.Logger
var AccessLog *logrus.Logger
var errorLogFile = "./log/error.log"
var accessLogFile = "./log/access.log"

func init() {
	os.MkdirAll("./log", os.ModePerm)
	initErrorLog()
	initAccessLog()
}

func initErrorLog() {
	ErrorLog = logrus.New()
	ErrorLog.SetFormatter(&logrus.JSONFormatter{})
	file, err := os.OpenFile(errorLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}
	ErrorLog.SetOutput(file)
}

func initAccessLog() {
	AccessLog = logrus.New()
	AccessLog.SetFormatter(&logrus.JSONFormatter{})
	file, err := os.OpenFile(accessLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}
	AccessLog.SetOutput(file)
}
