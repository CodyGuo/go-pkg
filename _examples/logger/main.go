package main

import (
	"errors"
	"github.com/CodyGuo/go-pkg/logger"
	"net"
)

const (
	logSender = "example"
	apiSender = "api"
	requestID = "cb12e64d-86af-4447-b66c-40c26a2e14f2"
)

func init() {
	conf := logger.Config{
		Level:               "info",
		TimeFormat:          "2006/01/02T15:04:05.000",
		FilePath:            "./log/logger.log",
		AccessFilePath:      "./log/logger_access.log",
		MaxSize:             10,
		MaxAge:              7,
		MaxBackups:          5,
		Compress:            false,
		UTCTime:             false,
		EnableFile:          true,
		EnableConsole:       true,
		EnableAccessFile:    false,
		EnableAccessConsole: true,
	}
	if err := conf.Init(); err != nil {
		return
		// log.Fatal(err)
	}
}

type User struct {
	Name string `json:"name"`
}

func (u User) Run(e *logger.Event, level logger.Level, msg string) {
	e.Str("name", u.Name)
}

func main() {
	logger.Info("info log to file and console")
	logger.InfoToFile("info log to file")
	logger.InfoToConsole("info log to console")
	logger.Error("error log to file and console")

	accessLogger := logger.GetAccessLogger()
	accessLogger.Info("info access log to file and console")
	accessLogger.InfoToFile("info access log to file")
	accessLogger.InfoToConsole("info access log to console")
	accessLogger.Error("error access log to file and console")

	var user = User{Name: "CodyGuo"}
	logger.With("str", "str").
		With("int", 1).
		With("[]int", []int{1, 2, 3}).
		With("ip", net.ParseIP("192.168.56.101")).
		With("user", user).
		WithCaller().
		WithSender(logSender).
		WithRequestID(requestID).
		WithHook(user).
		WithHookFunc(withErr(errors.New("hook error"))).
		Info("")

	apiLog := logger.WithSender(apiSender).WithRequestID(requestID).WithCaller()
	apiLog.WithHook(user).
		WithHookFunc(withErr(errors.New("hook error"))).
		Info("")
}

func withErr(err error) func(e *logger.Event, level logger.Level, message string) {
	return func(e *logger.Event, level logger.Level, message string) {
		e.Err(err)
	}
}
