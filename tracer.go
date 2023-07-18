package conflogger

import (
	"context"
	"time"

	"github.com/go-courier/logr"
	b3prop "go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

func InstallNewPipeline(outputType OutputType, formatType FormatType) error {
	stdout := StdoutSpanExporter(formatType)
	opts := []sdktrace.TracerProviderOption{
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSyncer(WithSpanMapExporter(OutputFilter(outputType))(stdout)),
	}

	tp := sdktrace.NewTracerProvider(opts...)
	otel.SetTextMapPropagator(b3prop.New())
	otel.SetTracerProvider(tp)
	return nil
}

func NewContextAndLogger(ctx context.Context, name string) (context.Context, logr.Logger) {
	ctx, span := otel.Tracer(name).Start(ctx, name, trace.WithTimestamp(time.Now()))
	log := SpanLogger(span)
	return logr.WithLogger(ctx, log), log
}

func SpanOnlyFilter() SpanMapper {
	return func(data sdktrace.ReadOnlySpan) sdktrace.ReadOnlySpan {
		if data == nil {
			return nil
		}

		d := &tracetest.SpanStub{}
		d.SpanContext = data.SpanContext()
		d.Parent = data.Parent()
		d.SpanKind = data.SpanKind()
		d.Name = data.Name()
		d.StartTime = data.StartTime()
		d.EndTime = data.EndTime()
		d.Attributes = data.Attributes()
		d.Links = data.Links()
		d.Status.Code = data.Status().Code
		d.Status.Description = data.Status().Description
		d.DroppedAttributes = data.DroppedAttributes()
		d.DroppedEvents = data.DroppedEvents()
		d.DroppedLinks = data.DroppedLinks()
		d.ChildSpanCount = data.ChildSpanCount()
		d.Resource = data.Resource()
		d.InstrumentationLibrary = data.InstrumentationLibrary()
		return d.Snapshot()
	}
}
