package alog

import (
	"context"
	"fmt"

	"github.com/alphaframework/alpha/autil/ahttp/request"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	defaultLogRootDir = "/data/log"
	defaultLogLevel   = "info"
	defaultLogFormat  = "json"

	RequestIdKey = "request_id"
)

var (
	Logger *zap.Logger
	Sugar  *zap.SugaredLogger
)

// InitLogger init Logger and Sugar
// applicationName
// rootDirectory: log root directory
// level: log level (debug/info/warn/error/panic/fatal)
func InitLogger(applicationName, rootDirectory, level, format string) error {
	if applicationName == "" {
		return fmt.Errorf("applicationName is required")
	}

	if rootDirectory == "" {
		rootDirectory = defaultLogRootDir
	}
	directory := rootDirectory + "/" + applicationName
	logPath := directory + "/" + applicationName + ".log"

	if level == "" {
		level = defaultLogLevel
	}
	if format == "" {
		format = defaultLogFormat
	}

	output := lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    512, // MB
		MaxAge:     240, // day
		MaxBackups: 100,
		Compress:   true,
	}
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	var l zapcore.Level
	if err := l.Set(level); err != nil {
		return err
	}
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(l)
	var writes = []zapcore.WriteSyncer{zapcore.AddSync(&output)}
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	if format == "console" {
		encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
		zapcore.NewMapObjectEncoder()
	}
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(writes...),
		atomicLevel,
	)

	caller := zap.AddCaller()
	development := zap.Development()

	field := zap.Fields(zap.String("application_name", applicationName))
	Logger = zap.New(core, caller, development, field)

	Logger.Info(`log.InitLogger successfully`)

	Sugar = Logger.Sugar()

	return nil
}

func CtxSugar(ctx context.Context) *zap.SugaredLogger {
	return Ctx(ctx).Sugar()
}

func Ctx(ctx context.Context) *zap.Logger {
	field := zap.Fields(zap.String(RequestIdKey, request.RequestIdValue(ctx)))
	return Logger.WithOptions(field)
}
