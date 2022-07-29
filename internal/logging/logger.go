package logging

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogLevel int

const (
	ErrorLevel = LogLevel(iota)
	DebugLeve
)

type correlationIdType int

const (
	requestIdKey correlationIdType = iota
	sessionIdKey
)

type Logger struct {
	*zap.Logger
}

func NewMockLogger() (logger *Logger) {
	config := zap.NewProductionEncoderConfig()

	config.EncodeTime = zapcore.ISO8601TimeEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(config)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel),
	)

	logger = &Logger{zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))}

	return logger
}

func NewZapLogger(logFileLvl, logConsoleLvl LogLevel) (logger *Logger) {
	config := zap.NewProductionEncoderConfig()

	config.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(config)

	consoleEncoder := zapcore.NewConsoleEncoder(config)

	logFile, _ := os.OpenFile("logs/log.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	writer := zapcore.AddSync(logFile)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, levelToZapTranslate(logFileLvl)),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), levelToZapTranslate(logConsoleLvl)),
	)

	logger = &Logger{zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))}

	return logger
}

func levelToZapTranslate(lvl LogLevel) zapcore.Level {
	switch lvl {
	case DebugLeve:
		return zap.DebugLevel
	case ErrorLevel:
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}

// WithRqId returns a context which knows its request ID
func WithRqId(ctx context.Context, rqId string) context.Context {
	return context.WithValue(ctx, requestIdKey, rqId)
}

// Logger returns a zap logger with as much context as possible
func (lg *Logger) WithCtx(ctx context.Context) *Logger {
	if ctx != nil {
		if ctxRqId, ok := ctx.Value(requestIdKey).(string); ok {

			return &Logger{lg.With(zap.String("rqId", ctxRqId))}
		}
	}
	return lg
}
