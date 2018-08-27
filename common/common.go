package common

import (
	"time"

	gocommon "github.com/liuhengloveyou/go-common"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type NilWriter struct{}

func (p *NilWriter) Write(b []byte) (n int, err error) { return 0, nil }

type Config struct {
	HTTPApiAddr string `toml:"http_api_addr"`
	DNSApiAddr  string `toml:"dns_api_addr"`
	CacheTTL    int64  `toml:"cache_ttl"`

	Mysql string `toml:"mysql"`

	GeoDB string `toml:"geodb"`

	LogDir   string `toml:"log_dir"`
	LogLevel string `toml:"log_level"`
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

type Rule struct {
	ID     int
	Domain string
	Line   string
	Area   string
	Zone   string
	Group  string
}

type Zone struct {
	ID   int
	Line string
	Area string
	Zone string
}

var (
	ServConfig Config
	Logger     *zap.Logger
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
		rotatelogs.WithMaxAge(60*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(writer),
		zap.DebugLevel)

	Logger = zap.New(core)
	Logger.WithOptions(zap.Development())
}

func ParseRRType(t uint16) string {

	switch t {
	case 1:
		return "A"
	case 5:
		return "CNAME"
	default:
		return ""
	}

	return ""
}
