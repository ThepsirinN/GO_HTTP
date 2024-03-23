package logger

import (
	"context"
	"fmt"
	"go_http_barko/config"
	"log"
	"os"
	"path/filepath"
	"time"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	ENV_LOCAL      = "local"
	ENV_DEVELOP    = "develop"
	ENV_PRODUCTION = "production"
)

const (
	TRACE_ID_KEY = "traceId"
	SPAN_ID_KEY  = "spanId"
)

var logger *zap.Logger

func InitLogger(cfg *config.Config) {
	logFile := initLogFile()
	encoderConfig := initEncoderConfig()

	var core zapcore.Core
	var lvl zapcore.Level

	if cfg.Log.Env == ENV_LOCAL {
		lvl = zap.InfoLevel
		coreLog := zapcore.NewCore(
			zapcore.NewJSONEncoder(*encoderConfig),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(logFile)),
			lvl,
		)
		core = zapcore.NewTee(coreLog)
	}

	logger = zap.New(core)

}

func initLogFile() *os.File {
	absPath, err := filepath.Abs("./log")
	if err != nil {
		log.Fatal("Error reading given path:", err)
	}
	logFile, err := os.OpenFile(absPath+"/logFile.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}

	return logFile
}

func initEncoderConfig() *zapcore.EncoderConfig {
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Local().Format("02 Jan 2006 15:04:05+07:00"))
	}

	return &zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     customTimeEncoder, // Using custom time format
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func Info(ctx context.Context, msg string, field ...zapcore.Field) {
	addTraceFromCtx(ctx, &field)
	logger.Info(msg, field...)
}

func Warn(ctx context.Context, msg string, field ...zapcore.Field) {
	addTraceFromCtx(ctx, &field)
	logger.Warn(msg, field...)
}

func Error(ctx context.Context, msg string, field ...zapcore.Field) {
	addTraceFromCtx(ctx, &field)
	logger.Error(msg, field...)
}

func Sync() {
	logger.Sync()
}

type traceInfo struct {
	traceId string
	spanId  string
}

func getTraceFromCtx(ctx context.Context) (isSpanContextValid bool, t traceInfo) {
	spanCtx := trace.SpanFromContext(ctx).SpanContext()
	fmt.Println(ctx)
	fmt.Printf("%+v\n", spanCtx)
	fmt.Printf("%+v\n", spanCtx.HasSpanID())
	fmt.Printf("%+v\n", spanCtx.HasTraceID())
	fmt.Println(spanCtx.IsValid())
	if spanCtx.IsValid() {
		t.traceId = spanCtx.TraceID().String()
		t.spanId = spanCtx.SpanID().String()
	}

	return false, t
}

func addTraceFromCtx(ctx context.Context, field *[]zapcore.Field) {
	isSpanContextValid, t := getTraceFromCtx(ctx)
	if isSpanContextValid {
		*field = append(*field, zap.Any(TRACE_ID_KEY, t.traceId), zap.Any(SPAN_ID_KEY, t.spanId))
	}
}
