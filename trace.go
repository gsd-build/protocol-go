package protocol

import (
	"crypto/rand"
	"fmt"
	"strings"
)

// TraceID extracts the trace-id component from a W3C traceparent header.
// Returns "" if the traceparent is empty or malformed.
// Format: version-traceid-parentid-traceflags (e.g. "00-abc...def-012...789-01")
func TraceID(traceparent string) string {
	if traceparent == "" {
		return ""
	}
	parts := strings.Split(traceparent, "-")
	if len(parts) < 4 || len(parts[1]) != 32 {
		return ""
	}
	return parts[1]
}

// NewTraceparent generates a new W3C traceparent header with random trace and span IDs.
// Format: 00-{32 hex trace-id}-{16 hex span-id}-01
func NewTraceparent() string {
	traceID := make([]byte, 16)
	spanID := make([]byte, 8)
	_, _ = rand.Read(traceID)
	_, _ = rand.Read(spanID)
	return fmt.Sprintf("00-%x-%x-01", traceID, spanID)
}

// ChildSpan creates a new traceparent with the same trace-id but a new span-id.
// Returns a new traceparent if the input is empty.
func ChildSpan(traceparent string) string {
	traceID := TraceID(traceparent)
	if traceID == "" {
		return NewTraceparent()
	}
	spanID := make([]byte, 8)
	_, _ = rand.Read(spanID)
	return fmt.Sprintf("00-%s-%x-01", traceID, spanID)
}
