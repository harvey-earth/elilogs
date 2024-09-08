package utils

import (
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
