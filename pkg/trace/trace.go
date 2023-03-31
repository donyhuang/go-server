package trace

import (
	"context"
	"github.com/hashicorp/go-uuid"
	"go.opentelemetry.io/otel/trace"
)

func IdFromContext(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		return spanCtx.TraceID().String()
	}
	traceId, _ := uuid.GenerateUUID()
	return traceId
}
