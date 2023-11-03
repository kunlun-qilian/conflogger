package conflogger

import (
	"context"
	"os"
	"runtime"

	"github.com/go-courier/metax"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/exp/slog"
)

type OutputType string

var (
	OutputAlways    OutputType = "Always"
	OutputOnFailure OutputType = "OnFailure"
	OutputNever     OutputType = "Never"
)

type FormatType string

var (
	FormatTEXT FormatType = "text"
	FormatJSON FormatType = "json"
)

func OutputFilter(outputType OutputType) SpanMapper {
	return func(data sdktrace.ReadOnlySpan) sdktrace.ReadOnlySpan {
		if outputType == OutputNever {
			return nil
		}
		if outputType == OutputOnFailure {
			if data.Status().Code == codes.Ok {
				return nil
			}
		}
		return data
	}
}

type SpanMapper = func(data sdktrace.ReadOnlySpan) sdktrace.ReadOnlySpan

func WithSpanMapExporter(mappers ...SpanMapper) func(spanExporter sdktrace.SpanExporter) sdktrace.SpanExporter {
	return func(spanExporter sdktrace.SpanExporter) sdktrace.SpanExporter {
		return &spanMapExporter{
			mappers:      mappers,
			SpanExporter: spanExporter,
		}
	}
}

type spanMapExporter struct {
	mappers []SpanMapper
	sdktrace.SpanExporter
}

func (e *spanMapExporter) ExportSpans(ctx context.Context, spanData []sdktrace.ReadOnlySpan) error {
	finalSpanSnapshot := make([]sdktrace.ReadOnlySpan, 0)

	mappers := e.mappers

	for i := range spanData {
		data := spanData[i]

		for _, m := range mappers {
			data = m(data)
		}

		if data != nil {
			finalSpanSnapshot = append(finalSpanSnapshot, data)
		}
	}

	if len(finalSpanSnapshot) == 0 {
		return nil
	}

	return e.SpanExporter.ExportSpans(ctx, finalSpanSnapshot)
}

func StdoutSpanExporter(formatType FormatType) sdktrace.SpanExporter {
	if formatType == FormatJSON {
		return &stdoutSpanExporter{formatter: slog.NewJSONHandler(os.Stdout, nil)}
	}
	return &stdoutSpanExporter{formatter: slog.NewTextHandler(os.Stdout, nil)}
}

type stdoutSpanExporter struct {
	formatter slog.Handler
}

func (e *stdoutSpanExporter) Shutdown(ctx context.Context) error {
	return nil
}

// ExportSpan writes a SpanSnapshot in json format to stdout.
func (e *stdoutSpanExporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {

	var pcs [1]uintptr
	// skip [runtime.Callers, this function, this function's caller]
	runtime.Callers(3, pcs[:])

	for i := range spans {
		data := spans[i]

		for _, event := range data.Events() {
			if event.Name == "" || event.Name[0] != '@' {
				continue
			}

			var lv slog.Level
			if err := lv.UnmarshalText([]byte(event.Name[1:])); err != nil {
				continue
			}

			// 使用slog，设置 slog 的时间
			log := slog.New(e.formatter)

			record := slog.NewRecord(event.Time, lv, "", pcs[0])

			for _, kv := range event.Attributes {
				k := string(kv.Key)

				switch k {
				case "msg":
					record.Message = kv.Value.AsString()
				default:
					log = log.With(slog.Any(k, kv.Value.AsInterface()))
				}
			}

			for _, kv := range data.Attributes() {
				k := string(kv.Key)
				if k == "meta" {
					meta := metax.ParseMeta(kv.Value.AsString())
					for k := range meta {
						log = log.With(slog.Any(k, meta[k]))
					}
					continue
				}
				log = log.With(slog.Any(k, kv.Value.AsInterface()))
			}

			log = log.With(slog.Any("span", data.Name()))

			if data.SpanContext().HasTraceID() {
				log = log.With(slog.Any("traceID", data.SpanContext().TraceID()))
			}

			if data.SpanContext().HasSpanID() {
				log = log.With(slog.Any("spanID", data.SpanContext().SpanID()))
			}

			if data.Parent().IsValid() {
				log = log.With(slog.Any("parentSpanID", data.Parent().SpanID()))
			}

			_ = log.Handler().Handle(ctx, record)
		}
	}

	return nil
}
