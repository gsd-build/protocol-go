package protocol

import "strings"

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
