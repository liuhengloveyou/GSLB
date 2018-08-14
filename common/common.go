package common

import (
	"strings"
	"time"

	gocommon "github.com/liuhengloveyou/go-common"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
)

type NilWriter struct{}

func (p *NilWriter) Write(b []byte) (n int, err error) { return 0, nil }

type Config struct {
	HTTPApiAddr string `toml:"http_api_addr"`
	DNSApiAddr  string `toml:"dns_api_addr"`
	Mysql       string `toml:"mysql"`
	LogDir      string `toml:"log_dir"`
	LogLevel    string `toml:"log_level"`
}

type RR struct {
	ID     int
	Domain string
	Ttl    uint32
	Type   uint16
	Class  uint16
	Data   string
	Group  string
}

var (
	ServConfig Config
)

func init() {
	if e := gocommon.LoadTomlConfig("./app.conf.toml", &ServConfig); e != nil {
		panic(e)
	}

	initLogger()
}

func initLogger() {
	writer, _ := rotatelogs.New(
		ServConfig.LogDir+"app.log.%Y%m%d%H%M",
		rotatelogs.WithLinkName(ServConfig.LogDir+"app.log"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)

	//log.SetOutput(&NilWriter{})
	log.SetFormatter(&log.TextFormatter{})

	log.AddHook(lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer,
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
	}, &log.TextFormatter{}))

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
