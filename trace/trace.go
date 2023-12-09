package trace

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func TraceName() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return fmt.Sprintf("%s:%d %s\n", frame.File, frame.Line, frame.Function)
}

func TraceFunction(ctx context.Context, tracer trace.Tracer) (context.Context, trace.Span) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	functionNameSplit := strings.Split(frame.Func.Name(), "/")
	functionName := functionNameSplit[len(functionNameSplit)-1]
	packageName := strings.Join(append(functionNameSplit[:len(functionNameSplit)-1], strings.Split(functionName, ".")[0]), "/")

	ctx, span := tracer.Start(ctx, functionName)

	fileDetails := fmt.Sprintf("%s#%d", frame.File, frame.Line)
	span.SetAttributes(
		attribute.String("function", functionName),
		attribute.String("function.package", packageName),
		attribute.String("function.file", fileDetails),
	)

	return ctx, span

}
