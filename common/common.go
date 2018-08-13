package common

import (
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
)

type NilWriter struct{}

func (p *NilWriter) Write(b []byte) (n int, err error) { return 0, nil }

type Config struct {
	Listen   string `json:"listen"`
	LogDir   string `json:"log_dir"`
	LogLevel string `json:"log_level"`
	Mysql    string `json:"mysql"`
}

var (
	ServConfig Config
)

func InitLogger() {
	writer, _ := rotatelogs.New(
		ServConfig.LogDir+"app.log.%Y%m%d%H%M",
		rotatelogs.WithLinkName(ServConfig.LogDir+"app.log"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)

	lfs := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer,
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
	}, &log.TextFormatter{})

	log.AddHook(lfs)
	log.SetOutput(&NilWriter{})
	log.SetFormatter(&log.TextFormatter{})

	logLvl := strings.ToLower(ServConfig.LogLevel)
	switch logLvl {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.DebugLevel)
	}
}
