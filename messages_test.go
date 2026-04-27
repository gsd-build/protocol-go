package protocol

import (
	"encoding/json"
	"testing"
	"time"
)

func floatPtr(v float64) *float64 { return &v }
func int64Ptr(v int64) *int64     { return &v }

func TestEnvelopeRoundTrip(t *testing.T) {
	cases := []struct {
		name string
		msg  any
	}{
		{"task", &Task{
			Type:            MsgTypeTask,
			TaskID:          "11111111-1111-1111-1111-111111111111",
			SessionID:       "22222222-2222-2222-2222-222222222222",
			ChannelID:       "ch-1",
			Prompt:          "hello",
			Engine:          "pi",
			Model:           "claude-opus-4-6[1m]",
			Effort:          "max",
			PermissionMode:  "acceptEdits",
			CWD:             "/tmp/project",
			ClaudeSessionID: "claude-abc-123",
			RequestID:       "33333333-3333-3333-3333-333333333333",
			Traceparent:     "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
		}},
		{"stream", &Stream{
			Type:           MsgTypeStream,
			SessionID:      "22222222-2222-2222-2222-222222222222",
			ChannelID:      "ch-1",
			SequenceNumber: 42,
			Event:          json.RawMessage(`{"delta":{"text":"hi"}}`),
			RequestID:      "33333333-3333-3333-3333-333333333333",
			Traceparent:    "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
		}},
		{"hello", &Hello{
			Type:          MsgTypeHello,
			MachineID:     "mach-id",
			DaemonVersion: "0.1.0",
			OS:            "darwin",
			Arch:          "arm64",
		}},
		{"hello capabilities", &Hello{
			Type:          MsgTypeHello,
			MachineID:     "mach-id",
			DaemonVersion: "0.2.0",
			OS:            "darwin",
			Arch:          "arm64",
			Capabilities: &HelloCapabilities{
				Stop: true,
			},
			ActiveTasks: []string{"task-a", "task-b"},
		}},
		{"welcome", &Welcome{
			Type:                MsgTypeWelcome,
			LatestDaemonVersion: "0.2.1",
		}},
		{"taskComplete", &TaskComplete{
			Type:            MsgTypeTaskComplete,
			TaskID:          "11111111-1111-1111-1111-111111111111",
			SessionID:       "22222222-2222-2222-2222-222222222222",
			ChannelID:       "ch-1",
			ClaudeSessionID: "claude-abc",
			InputTokens:     100,
			OutputTokens:    50,
			CostUSD:         "0.0150",
			DurationMs:      1234,
			RequestID:       "33333333-3333-3333-3333-333333333333",
			Traceparent:     "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
		}},
		{"taskStarted", &TaskStarted{
			Type:        MsgTypeTaskStarted,
			TaskID:      "11111111-1111-1111-1111-111111111111",
			SessionID:   "22222222-2222-2222-2222-222222222222",
			ChannelID:   "ch-1",
			StartedAt:   "2026-04-13T12:00:00Z",
			RequestID:   "33333333-3333-3333-3333-333333333333",
			Traceparent: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
		}},
		{"taskError", &TaskError{
			Type:        MsgTypeTaskError,
			TaskID:      "11111111-1111-1111-1111-111111111111",
			SessionID:   "22222222-2222-2222-2222-222222222222",
			ChannelID:   "ch-1",
			Error:       "boom",
			RequestID:   "33333333-3333-3333-3333-333333333333",
			Traceparent: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
		}},
		{"taskCancelled", &TaskCancelled{
			Type:        MsgTypeTaskCancelled,
			TaskID:      "11111111-1111-1111-1111-111111111111",
			SessionID:   "22222222-2222-2222-2222-222222222222",
			ChannelID:   "ch-1",
			RequestID:   "33333333-3333-3333-3333-333333333333",
			Traceparent: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
		}},
		{"question", &Question{
			Type:        MsgTypeQuestion,
			SessionID:   "22222222-2222-2222-2222-222222222222",
			ChannelID:   "ch-1",
			RequestID:   "33333333-3333-3333-3333-333333333333",
			Question:    "Which library should we use?",
			Header:      "Library",
			MultiSelect: true,
			Options: []QuestionOption{
				{
					Label:       "date-fns",
					Description: "Tree-shakable, functional API, smaller bundle",
					Preview:     "import { format } from 'date-fns'",
				},
				{
					Label:       "dayjs",
					Description: "Moment-style API, plugin ecosystem",
				},
			},
		}},
		{"questionResponse", &QuestionResponse{
			Type:      MsgTypeQuestionResponse,
			SessionID: "22222222-2222-2222-2222-222222222222",
			ChannelID: "ch-1",
			RequestID: "33333333-3333-3333-3333-333333333333",
			Answer:    `["date-fns","custom note"]`,
		}},
		{"compact request", &CompactRequest{
			Type:         MsgTypeCompactRequest,
			SessionID:    "session_123",
			ChannelID:    "channel_123",
			RequestID:    "compact_123",
			Instructions: "preserve auth state and exact file paths",
		}},
		{"context stats request", &ContextStatsRequest{
			Type:      MsgTypeContextStatsRequest,
			SessionID: "session_123",
			ChannelID: "channel_123",
			RequestID: "stats_123",
		}},
		{"context stats", &ContextStats{
			Type:                 MsgTypeContextStats,
			SessionID:            "session_123",
			ChannelID:            "channel_123",
			RequestID:            "stats_123",
			Tokens:               int64Ptr(270000),
			ContextWindow:        1000000,
			Percent:              floatPtr(27.0),
			ReserveTokens:        16384,
			KeepRecentTokens:     20000,
			AutoThresholdPercent: 98.3616,
			Source:               "pi",
			ObservedAt:           time.Date(2026, 4, 27, 12, 0, 0, 0, time.UTC),
		}},
		{"compact status completed", &CompactStatus{
			Type:                 MsgTypeCompactStatus,
			SessionID:            "session_123",
			ChannelID:            "channel_123",
			RequestID:            "compact_123",
			Status:               CompactStatusCompleted,
			Reason:               CompactReasonManual,
			Instructions:         "preserve auth state and exact file paths",
			TokensBefore:         int64Ptr(8951),
			TokensAfter:          int64Ptr(7712),
			ContextWindow:        1000000,
			ReserveTokens:        16384,
			KeepRecentTokens:     20000,
			AutoThresholdPercent: 98.3616,
			Summary:              "The session is working on Pi context compaction.",
			FirstKeptEntryID:     "entry_42",
			Source:               "pi",
			ObservedAt:           time.Date(2026, 4, 27, 12, 1, 0, 0, time.UTC),
		}},
		{"previewOpen", &PreviewOpen{
			Type:       MsgTypePreviewOpen,
			RequestID:  "req-1",
			PreviewID:  "preview_123",
			SessionID:  "session_123",
			ChannelID:  "channel_123",
			MachineID:  "machine_123",
			TargetHost: "127.0.0.1",
			TargetPort: 3000,
			ExpiresAt:  "2026-04-27T20:00:00Z",
		}},
		{"previewOpenResult", &PreviewOpenResult{
			Type:      MsgTypePreviewOpenResult,
			RequestID: "req-1",
			PreviewID: "preview_123",
			OK:        true,
		}},
		{"previewHttpRequest", &PreviewHTTPRequest{
			Type:      MsgTypePreviewHTTPRequest,
			RequestID: "req-2",
			StreamID:  "stream_1",
			PreviewID: "preview_123",
			Method:    "POST",
			Path:      "/api/action",
			Headers: map[string][]string{
				"host":              {"preview_123.preview.gsd.build"},
				"x-forwarded-proto": {"https"},
			},
		}},
		{"previewHttpResponseHead", &PreviewHTTPResponseHead{
			Type:       MsgTypePreviewHTTPResponseHead,
			RequestID:  "req-2",
			StreamID:   "stream_1",
			PreviewID:  "preview_123",
			StatusCode: 200,
			Headers: map[string][]string{
				"content-type": {"text/html; charset=utf-8"},
			},
		}},
		{"previewStreamChunk", &PreviewStreamChunk{
			Type:       MsgTypePreviewStreamChunk,
			StreamID:   "stream_1",
			Sequence:   1,
			BodyBase64: "aGVsbG8=",
			Final:      false,
		}},
		{"previewStreamCancel", &PreviewStreamCancel{
			Type:     MsgTypePreviewStreamCancel,
			StreamID: "stream_1",
			Reason:   "browser_abort",
		}},
		{"previewWebSocketOpen", &PreviewWebSocketOpen{
			Type:      MsgTypePreviewWebSocketOpen,
			StreamID:  "ws_1",
			PreviewID: "preview_123",
			Path:      "/_next/webpack-hmr",
			Headers:   map[string][]string{},
			Protocols: []string{"vite-hmr"},
		}},
		{"previewWebSocketOpenResult", &PreviewWebSocketOpenResult{
			Type:      MsgTypePreviewWebSocketOpenResult,
			StreamID:  "ws_1",
			PreviewID: "preview_123",
			OK:        true,
			Protocol:  "vite-hmr",
		}},
		{"previewWebSocketData", &PreviewWebSocketData{
			Type:       MsgTypePreviewWebSocketData,
			StreamID:   "ws_1",
			Sequence:   1,
			IsBinary:   false,
			BodyBase64: "eyJ0eXBlIjoicGluZyJ9",
		}},
		{"previewWebSocketClose", &PreviewWebSocketClose{
			Type:     MsgTypePreviewWebSocketClose,
			StreamID: "ws_1",
			Code:     1000,
			Reason:   "normal",
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.msg)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}

			env, err := ParseEnvelope(data)
			if err != nil {
				t.Fatalf("parse envelope: %v", err)
			}

			// Round-trip should preserve the original JSON
			reMarshaled, err := json.Marshal(env.Payload)
			if err != nil {
				t.Fatalf("re-marshal: %v", err)
			}

			// Parse both into maps and compare, to ignore field ordering
			var original, final map[string]any
			if err := json.Unmarshal(data, &original); err != nil {
				t.Fatalf("unmarshal original: %v", err)
			}
			if err := json.Unmarshal(reMarshaled, &final); err != nil {
				t.Fatalf("unmarshal round trip: %v", err)
			}

			if !jsonEqual(original, final) {
				t.Errorf("payload mismatch after round trip: want %v, got %v", original, final)
			}
		})
	}
}

func jsonEqual(a, b any) bool {
	ja, _ := json.Marshal(a)
	jb, _ := json.Marshal(b)
	return string(ja) == string(jb)
}

func TestHelloPreviewCapabilitiesRoundTrip(t *testing.T) {
	msg := &Hello{
		Type:          MsgTypeHello,
		MachineID:     "machine_123",
		DaemonVersion: "0.5.0",
		OS:            "darwin",
		Arch:          "arm64",
		Capabilities: &HelloCapabilities{
			Stop:                      true,
			PreviewTunnel:             true,
			PreviewMaxFrameBytes:      1048576,
			PreviewChunkBytes:         196608,
			PreviewWebSocketProtocols: true,
		},
	}
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	env, err := ParseEnvelope(data)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	got := env.Payload.(*Hello)
	if got.Capabilities == nil || !got.Capabilities.PreviewTunnel {
		t.Fatalf("preview capability missing after round trip: %#v", got.Capabilities)
	}
}

func TestParseEnvelopeRejectsUnknownType(t *testing.T) {
	_, err := ParseEnvelope([]byte(`{"type":"bogus"}`))
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestHelloOmitsEmptyActiveTasks(t *testing.T) {
	h := &Hello{
		Type:          MsgTypeHello,
		MachineID:     "m-1",
		DaemonVersion: "0.2.0",
		OS:            "darwin",
		Arch:          "arm64",
	}

	data, err := json.Marshal(h)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, exists := raw["activeTasks"]; exists {
		t.Error("activeTasks should be omitted when empty")
	}
}

func TestWelcomeOmitsEmptyVersion(t *testing.T) {
	w := &Welcome{
		Type: MsgTypeWelcome,
	}

	data, err := json.Marshal(w)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, exists := raw["latestDaemonVersion"]; exists {
		t.Error("latestDaemonVersion should be omitted when empty")
	}
}

func TestTraceparentOmittedWhenEmpty(t *testing.T) {
	task := &Task{
		Type:           MsgTypeTask,
		TaskID:         "11111111-1111-1111-1111-111111111111",
		SessionID:      "22222222-2222-2222-2222-222222222222",
		ChannelID:      "ch-1",
		Prompt:         "hello",
		Model:          "claude-opus-4-6[1m]",
		Effort:         "max",
		PermissionMode: "acceptEdits",
		CWD:            "/tmp",
	}

	data, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, exists := raw["traceparent"]; exists {
		t.Error("traceparent should be omitted when empty")
	}
}

func TestTraceID(t *testing.T) {
	cases := []struct {
		traceparent string
		want        string
	}{
		{"00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01", "4bf92f3577b34da6a3ce929d0e0e4736"},
		{"", ""},
		{"invalid", ""},
		{"00-short-00f067aa0ba902b7-01", ""},
	}

	for _, tc := range cases {
		got := TraceID(tc.traceparent)
		if got != tc.want {
			t.Errorf("TraceID(%q) = %q, want %q", tc.traceparent, got, tc.want)
		}
	}
}

func TestMachineStatusRoundTrip(t *testing.T) {
	msg := &MachineStatus{
		Type:          MsgTypeMachineStatus,
		MachineID:     "machine-123",
		State:         "reconnecting",
		PreviousState: "online",
		Reason:        "lease_expired",
		OccurredAt:    "2026-04-15T12:00:00Z",
	}

	buf, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	env, err := ParseEnvelope(buf)
	if err != nil {
		t.Fatalf("parse envelope: %v", err)
	}

	got, ok := env.Payload.(*MachineStatus)
	if !ok {
		t.Fatalf("payload type = %T", env.Payload)
	}

	if got.State != "reconnecting" || got.PreviousState != "online" || got.Reason != "lease_expired" {
		t.Fatalf("unexpected payload: %+v", got)
	}
}
