package utils

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Elilogger is the logging instance
var Elilogger *zap.SugaredLogger

// InitLogger initializes the zap logger for use
func InitLogger() {
	var verbosity zapcore.Level
	switch lvl := viper.GetString("logLevel"); lvl {
	case "warn":
		verbosity = zap.WarnLevel
	case "info":
		verbosity = zap.InfoLevel
	case "debug":
		verbosity = zap.DebugLevel
	case "quiet":
		verbosity = zap.FatalLevel
	}

	cfg := zap.NewProductionEncoderConfig()
	cfg.TimeKey = "timestamp"
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(verbosity),
		DisableStacktrace: true,
		Encoding:          "console",
		OutputPaths: []string{
			"stdout",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		EncoderConfig: cfg,
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	Elilogger = logger.Sugar()
}

// LogRequest takes an esapi response and logs it
func LogRequest(r []byte) {
	Debug(string(r))
}

func cleanup(errString string, err error) string {
	errString = strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, errString)
	return fmt.Sprintf(strings.Replace(errString, "\n ", "", 1), err.Error())
}

// Debug returns a debug message to the logger
func Debug(debugString string) {
	Elilogger.Debug(debugString)
}

// Error returns an error message to the logger and terminates
func Error(errString string, err error) error {
	errString = cleanup(errString, err)
	if err != nil {
		Elilogger.Error(errString)
		os.Exit(1)
	}
	return nil
}

// Fatal returns a fatal msg to the logger and terminates
func Fatal(errString string, err error) {
	errString = cleanup(errString, err)
	if err != nil {
		Elilogger.Fatal(errString)
	}
}

// Info returns an info msg to the logger
// is info or debug
func Info(infoString string) {
	Elilogger.Info(infoString)
}

// Warn returns a warning msg to the logger
func Warn(warnString string) {
	Elilogger.Warn(warnString)
}

