package conflogger

import (
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type Log struct {
	ReportCaller bool
	Name         string
	Level        string `env:""`
	Format       string
	init         bool
}

func (log *Log) SetDefaults() {
	log.ReportCaller = true

	if log.Level == "" {
		log.Level = "DEBUG"
	}

	if log.Format == "" {
		log.Format = "json"
	}
}

func (log *Log) Init() {
	if !log.init {
		log.Create()
		log.init = true
	}
}

func (log *Log) Create() {
	log.SetDefaults()
	if log.Format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{
			//PrettyPrint:      true,
			CallerPrettyfier: CallerPrettyfier,
		})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:      true,
			CallerPrettyfier: CallerPrettyfier,
		})
	}

	logrus.SetLevel(getLogLevel(log.Level))
	logrus.SetReportCaller(log.ReportCaller)
	logrus.SetOutput(os.Stdout)
}

func getLogLevel(l string) logrus.Level {
	level, err := logrus.ParseLevel(strings.ToLower(l))
	if err == nil {
		return level
	}
	return logrus.InfoLevel
}

func CallerPrettyfier(f *runtime.Frame) (function string, file string) {
	return f.Function + " line:" + strconv.FormatInt(int64(f.Line), 10), ""
}
