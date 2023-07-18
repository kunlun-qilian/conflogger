package conflogger

import (
	"context"
	"errors"
	"github.com/go-courier/logr"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var logger = Log{
	Level:  "DEBUG",
	Format: "json",
}

func init() {
	logger.SetDefaults()
	logger.Init()
}

func TestLog(t *testing.T) {
	ctx := context.Background()
	doLog(ctx)

}

func doLog(ctx context.Context) {
	tracer := otel.Tracer("")

	ctx, span := tracer.Start(ctx, "op", trace.WithTimestamp(time.Now()))
	defer func() {
		span.End(trace.WithTimestamp(time.Now()))
	}()

	ctx = logr.WithLogger(ctx, SpanLogger(span))

	someActionWithSpan(ctx)

	otherActions(ctx)
}

func someActionWithSpan(ctx context.Context) {
	_, log := logr.Start(ctx, "SomeActionWithSpan")
	defer log.End()

	log.Info("info")
	log.Debug("debug")
	log.Warn(errors.New("warn"))
}

func otherActions(ctx context.Context) {
	log := logr.FromContext(ctx)

	log.WithValues("test2", 2).Info("test")
	log.Error(errors.New(""))
}
