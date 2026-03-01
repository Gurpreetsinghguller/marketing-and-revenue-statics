package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var appLogger = logrus.New()

func init() {
	appLogger.SetOutput(os.Stdout)
	appLogger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		TimestampFormat:  time.RFC3339,
		PadLevelText:     true,
		QuoteEmptyFields: true,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			return "", fmt.Sprintf("%s:%d", filepath.Base(frame.File), frame.Line)
		},
	})
	appLogger.SetReportCaller(true)
	appLogger.SetLevel(logrus.InfoLevel)
}

func Configure(level string) {
	appLogger.SetLevel(parseLevel(level))
}

func Get() *logrus.Logger {
	return appLogger
}

func parseLevel(level string) logrus.Level {
	parsed, err := logrus.ParseLevel(strings.ToLower(strings.TrimSpace(level)))
	if err != nil {
		return logrus.InfoLevel
	}
	return parsed
}
