package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log Instance
var Log *zap.SugaredLogger

func init() {
	// Init the logger when the package is used
	initLogger()
}

// InitLogger Init the logger
func initLogger() {
	var logLevel zapcore.Level
	logLevelStr := strings.ToUpper(os.Getenv("LOG_LEVEL"))
	switch logLevelStr {
	case "0", "INFO":
		logLevel = zapcore.InfoLevel
	case "1", "WARN":
		logLevel = zapcore.WarnLevel
	case "2", "ERROR":
		logLevel = zapcore.ErrorLevel
	case "3", "DPANIC":
		logLevel = zapcore.DPanicLevel
	case "4", "PANIC":
		logLevel = zapcore.PanicLevel
	case "5", "FATAL":
		logLevel = zapcore.FatalLevel
	case "9", "DEBUG":
		logLevel = zapcore.DebugLevel
	default:
		logLevel = zapcore.ErrorLevel
	}

	// First, define our level-handling logic.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return (lvl < zapcore.ErrorLevel) && (lvl >= logLevel)
	})

	// Console output
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Console encoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	// From a zapcore.Core, it's easy to construct a Logger.
	Log = zap.New(core).Sugar()

	defer Log.Sync()
	Log.Infof("successfully initiated logging, level = %s", logLevel.String())
}
