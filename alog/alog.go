package alog

import (
	"context"
	"fmt"

	"github.com/alphaframework/alpha/autil"
	"github.com/alphaframework/alpha/autil/ahttp/request"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	defaultLogLevel  = "info"
	defaultLogFormat = "console"

	RequestIdKey = "request_id"
)

var (
	Logger *zap.Logger
	Sugar  *zap.SugaredLogger
)

// InitLogger init Logger and Sugar
// applicationName
// directory: log directory
// level: log level (debug/info/warn/error/panic/fatal)
// format: log format (console/json)
func InitLogger(applicationName, directory, level, format string) error {
	if applicationName == "" {
		return fmt.Errorf("applicationName is required")
	}
	formatList := []string{"", "console", "json"}
	if !autil.In(format, formatList) {
		return fmt.Errorf("log format: %s does not validate as in %#v", format, formatList)
	}

	if directory == "" {
		return fmt.Errorf("directory is required")
	}

	if level == "" {
		level = defaultLogLevel
	}
	if format == "" {
		format = defaultLogFormat
	}

	var l zapcore.Level
	if err := l.Set(level); err != nil {
		return err
	}

	getWriter := func(level string) *lumberjack.Logger {
		logPath := directory + "/" + applicationName + "." + level + ".log"
		return &lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    512, // MB
			MaxAge:     240, // day
			MaxBackups: 100,
			Compress:   true,
		}
	}

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	if format == "json" {
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	var cores []zapcore.Core
	for l <= zapcore.FatalLevel {
		var writers = []zapcore.WriteSyncer{zapcore.AddSync(getWriter(l.String()))}
		cores = append(cores, zapcore.NewCore(
			encoder,
			zapcore.NewMultiWriteSyncer(writers...),
			zap.NewAtomicLevelAt(l),
		))
		l++
	}
	core := zapcore.NewTee(cores...)

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
