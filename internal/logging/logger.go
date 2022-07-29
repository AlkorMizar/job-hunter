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

type correlationIDType int

const (
	requestIDKey correlationIDType = iota
)

type Logger struct {
	*zap.Logger
}

const filePerm = 0o644

func NewDefaultLogger(lvl LogLevel) (logger *Logger) {
	config := zap.NewProductionEncoderConfig()

	config.EncodeTime = zapcore.ISO8601TimeEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(config)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), levelToZapTranslate(lvl)),
	)

	logger = &Logger{zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))}

	return logger
}

func NewZapLogger(logFileLvl, logConsoleLvl LogLevel, logFilePath string) (logger *Logger) {
	config := zap.NewProductionEncoderConfig()

	config.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(config)

	consoleEncoder := zapcore.NewConsoleEncoder(config)

	logFile, _ := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, filePerm)

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

// WithRqID returns a context which knows its request ID
func WithRqID(ctx context.Context, rqID string) context.Context {
	return context.WithValue(ctx, requestIDKey, rqID)
}

// Logger returns a zap logger with as much context as possible
func (lg *Logger) WithCtx(ctx context.Context) *Logger {
	if ctx != nil {
		if ctxRqID, ok := ctx.Value(requestIDKey).(string); ok {
			return &Logger{lg.With(zap.String("rqID", ctxRqID))}
		}
	}

	return lg
}
