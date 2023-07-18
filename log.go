package conflogger

import (
	"github.com/go-courier/logr"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strconv"
	"strings"
)

var l = logr.DebugLevel

func SetLevel(lvl logr.Level) {
	l = lvl
}

type Log struct {
	ReportCaller bool
	Level        string     `env:""`
	Output       OutputType `env:""`
	Format       FormatType
	init         bool
}

func (l *Log) SetDefaults() {

	if l.Output == "" {
		l.Output = OutputAlways
	}

	if l.Format == "" {
		l.Format = FormatJSON
	}
}

func (l *Log) Init() {
	if !l.init {
		if l.Format == "json" {
			logrus.SetFormatter(&logrus.JSONFormatter{
				CallerPrettyfier: CallerPrettyfier,
			})
		} else {
			logrus.SetFormatter(&logrus.TextFormatter{
				ForceColors:      true,
				CallerPrettyfier: CallerPrettyfier,
			})
		}

		logursLevel, logrLevel := getLogLevel(l.Level)
		SetLevel(logrLevel)
		logrus.SetLevel(logursLevel)
		logrus.SetReportCaller(l.ReportCaller)
		logrus.SetOutput(os.Stdout)

		if err := InstallNewPipeline(l.Output, l.Format); err != nil {
			panic(err)
		}
		l.init = true
	}
}

func getLogLevel(l string) (logrus.Level, logr.Level) {
	level, err := logrus.ParseLevel(strings.ToLower(l))
	if err != nil {
		return logrus.DebugLevel, logr.DebugLevel
	}
	logrLevel, logrLevelErr := logr.ParseLevel(l)
	if logrLevelErr != nil {
		return logrus.DebugLevel, logr.DebugLevel
	}
	return level, logrLevel
}

func CallerPrettyfier(f *runtime.Frame) (function string, file string) {
	return f.Function + " line:" + strconv.FormatInt(int64(f.Line), 10), ""
}
