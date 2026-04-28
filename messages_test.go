package protocol

import (
	"encoding/json"
	"testing"
	"time"
)

func floatPtr(v float64) *float64 { return &v }
func int64Ptr(v int64) *int64     { return &v }
func intPtr(v int) *int           { return &v }

func TestEnvelopeRoundTrip(t *testing.T) {
	cases := []struct {
		name string
		msg  any
	}{
		{"task", &Task{
			Type:               MsgTypeTask,
			TaskID:             "11111111-1111-1111-1111-111111111111",
			SessionID:          "22222222-2222-2222-2222-222222222222",
			ChannelID:          "ch-1",
			Prompt:             "hello",
			Engine:             "pi",
			Model:              "claude-opus-4-6[1m]",
			Effort:             "max",
			PermissionMode:     "acceptEdits",
			CWD:                "/tmp/project",
			ClaudeSessionID:    "claude-abc-123",
			RequestID:          "33333333-3333-3333-3333-333333333333",
			Traceparent:        "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
			CustomInstructions: "Always talk like a pirate.",
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
				Stop:     true,
				Terminal: true,
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
		{"browseDir paginated", &BrowseDir{
			Type:      MsgTypeBrowseDir,
			RequestID: "browse-1",
			ChannelID: "chan-1",
			MachineID: "machine-1",
			Path:      "/tmp/project",
			Limit:     200,
			Cursor:    "200",
		}},
		{"browseDirResult paginated", &BrowseDirResult{
			Type:       MsgTypeBrowseDirResult,
			RequestID:  "browse-1",
			ChannelID:  "chan-1",
			OK:         true,
			HasMore:    true,
			NextCursor: "400",
			Entries: []BrowseEntry{
				{
					Name:        "src",
					Path:        "/tmp/project/src",
					IsDirectory: true,
					Size:        64,
					ModifiedAt:  "2026-04-27T18:00:00Z",
				},
			},
		}},
		{"listSkills", &ListSkills{
			Type:      MsgTypeListSkills,
			RequestID: "skills-1",
			ChannelID: "chan-1",
			MachineID: "machine-1",
			CWD:       "/tmp/project",
		}},
		{"listSkillsResult", &ListSkillsResult{
			Type:      MsgTypeListSkillsResult,
			RequestID: "skills-1",
			ChannelID: "chan-1",
			OK:        true,
			Skills: []Skill{
				{
					Name:        "debug-like-expert",
					Description: "Deep analysis debugging workflow",
					Path:        "/Users/me/.claude/skills/debug-like-expert/SKILL.md",
					Scope:       "home",
				},
			},
		}},
		{"terminalOpen", &TerminalOpen{
			Type:      MsgTypeTerminalOpen,
			RequestID: "open-1",
			SessionID: "sess-1",
			ChannelID: "chan-1",
			Token:     "tok",
			Cols:      120,
			Rows:      32,
		}},
		{"terminalInput", &TerminalInput{
			Type:       MsgTypeTerminalInput,
			TerminalID: "term-1",
			ChannelID:  "chan-1",
			DataBase64: "YQ==",
		}},
		{"terminalOutput", &TerminalOutput{
			Type:       MsgTypeTerminalOutput,
			TerminalID: "term-1",
			SessionID:  "sess-1",
			ChannelID:  "chan-1",
			Seq:        7,
			DataBase64: "b2s=",
		}},
		{"terminalResize", &TerminalResize{
			Type:       MsgTypeTerminalResize,
			TerminalID: "term-1",
			ChannelID:  "chan-1",
			Cols:       100,
			Rows:       28,
		}},
		{"terminalClose", &TerminalClose{
			Type:       MsgTypeTerminalClose,
			TerminalID: "term-1",
			ChannelID:  "chan-1",
		}},
		{"terminalOpened", &TerminalOpened{
			Type:       MsgTypeTerminalOpened,
			RequestID:  "open-1",
			TerminalID: "term-1",
			SessionID:  "sess-1",
			ChannelID:  "chan-1",
			Shell:      "/bin/zsh",
			CWD:        "/tmp/project",
			StartedAt:  "2026-04-27T18:00:00Z",
		}},
		{"terminalSnapshot", &TerminalSnapshot{
			Type:       MsgTypeTerminalSnapshot,
			TerminalID: "term-1",
			SessionID:  "sess-1",
			ChannelID:  "chan-1",
			Seq:        8,
			DataBase64: "c25hcA==",
		}},
		{"terminalExit", &TerminalExit{
			Type:       MsgTypeTerminalExit,
			TerminalID: "term-1",
			SessionID:  "sess-1",
			ChannelID:  "chan-1",
			ExitCode:   intPtr(0),
			Reason:     "process_exit",
			EndedAt:    "2026-04-27T18:30:00Z",
		}},
		{"terminalError", &TerminalError{
			Type:       MsgTypeTerminalError,
			RequestID:  "open-1",
			TerminalID: "term-1",
			SessionID:  "sess-1",
			ChannelID:  "chan-1",
			Error:      "Unable to start shell",
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
		{"localServerDetected", &LocalServerDetected{
			Type:       MsgTypeLocalServerDetected,
			SessionID:  "session_123",
			ChannelID:  "channel_123",
			TaskID:     "task_123",
			ToolUseID:  "toolu_123",
			Host:       "127.0.0.1",
			Port:       5173,
			URL:        "http://127.0.0.1:5173/",
			Command:    "pnpm dev",
			Source:     "tool_output",
			DetectedAt: "2026-04-27T20:00:00Z",
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

func TestTerminalEnvelopeRejectsInvalidFieldTypes(t *testing.T) {
	cases := []struct {
		name string
		raw  string
	}{
		{
			name: "terminalOpen cols string",
			raw:  `{"type":"terminalOpen","requestId":"open-1","sessionId":"sess-1","channelId":"chan-1","cols":"120","rows":32}`,
		},
		{
			name: "terminalInput dataBase64 number",
			raw:  `{"type":"terminalInput","terminalId":"term-1","channelId":"chan-1","dataBase64":123}`,
		},
		{
			name: "terminalOutput seq string",
			raw:  `{"type":"terminalOutput","terminalId":"term-1","sessionId":"sess-1","channelId":"chan-1","seq":"7","dataBase64":"b2s="}`,
		},
		{
			name: "terminalResize rows string",
			raw:  `{"type":"terminalResize","terminalId":"term-1","channelId":"chan-1","cols":100,"rows":"28"}`,
		},
		{
			name: "terminalOpened startedAt object",
			raw:  `{"type":"terminalOpened","requestId":"open-1","terminalId":"term-1","sessionId":"sess-1","channelId":"chan-1","shell":"/bin/zsh","cwd":"/tmp/project","startedAt":{}}`,
		},
		{
			name: "terminalSnapshot seq object",
			raw:  `{"type":"terminalSnapshot","terminalId":"term-1","sessionId":"sess-1","channelId":"chan-1","seq":{},"dataBase64":"c25hcA=="}`,
		},
		{
			name: "terminalExit exitCode string",
			raw:  `{"type":"terminalExit","terminalId":"term-1","sessionId":"sess-1","channelId":"chan-1","exitCode":"0","reason":"process_exit","endedAt":"2026-04-27T18:30:00Z"}`,
		},
		{
			name: "terminalError error array",
			raw:  `{"type":"terminalError","requestId":"open-1","terminalId":"term-1","sessionId":"sess-1","channelId":"chan-1","error":[]}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := ParseEnvelope([]byte(tc.raw)); err == nil {
				t.Fatal("expected parse error")
			}
		})
	}
}

func TestLocalServerDetectedRejectsInvalidFieldTypes(t *testing.T) {
	cases := []struct {
		name string
		raw  string
	}{
		{
			name: "session id number",
			raw:  `{"type":"` + MsgTypeLocalServerDetected + `","sessionId":123,"channelId":"channel_123","host":"127.0.0.1","port":5173,"url":"http://127.0.0.1:5173/","detectedAt":"2026-04-27T20:00:00Z"}`,
		},
		{
			name: "port string",
			raw:  `{"type":"` + MsgTypeLocalServerDetected + `","sessionId":"session_123","channelId":"channel_123","host":"127.0.0.1","port":"5173","url":"http://127.0.0.1:5173/","detectedAt":"2026-04-27T20:00:00Z"}`,
		},
		{
			name: "source object",
			raw:  `{"type":"` + MsgTypeLocalServerDetected + `","sessionId":"session_123","channelId":"channel_123","host":"127.0.0.1","port":5173,"url":"http://127.0.0.1:5173/","source":{},"detectedAt":"2026-04-27T20:00:00Z"}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := ParseEnvelope([]byte(tc.raw)); err == nil {
				t.Fatal("expected parse error")
			}
		})
	}
}

func TestTerminalEnvelopeIgnoresUnknownFields(t *testing.T) {
	cases := []struct {
		name string
		raw  string
	}{
		{
			name: "terminalOpen",
			raw:  `{"type":"terminalOpen","requestId":"open-1","sessionId":"sess-1","channelId":"chan-1","token":"tok","cols":120,"rows":32,"unknown":"ok"}`,
		},
		{
			name: "terminalInput",
			raw:  `{"type":"terminalInput","terminalId":"term-1","channelId":"chan-1","dataBase64":"YQ==","unknown":"ok"}`,
		},
		{
			name: "terminalOutput",
			raw:  `{"type":"terminalOutput","terminalId":"term-1","sessionId":"sess-1","channelId":"chan-1","seq":7,"dataBase64":"b2s=","unknown":"ok"}`,
		},
		{
			name: "terminalResize",
			raw:  `{"type":"terminalResize","terminalId":"term-1","channelId":"chan-1","cols":100,"rows":28,"unknown":"ok"}`,
		},
		{
			name: "terminalClose",
			raw:  `{"type":"terminalClose","terminalId":"term-1","channelId":"chan-1","unknown":"ok"}`,
		},
		{
			name: "terminalOpened",
			raw:  `{"type":"terminalOpened","requestId":"open-1","terminalId":"term-1","sessionId":"sess-1","channelId":"chan-1","shell":"/bin/zsh","cwd":"/tmp/project","startedAt":"2026-04-27T18:00:00Z","unknown":"ok"}`,
		},
		{
			name: "terminalSnapshot",
			raw:  `{"type":"terminalSnapshot","terminalId":"term-1","sessionId":"sess-1","channelId":"chan-1","seq":8,"dataBase64":"c25hcA==","unknown":"ok"}`,
		},
		{
			name: "terminalExit",
			raw:  `{"type":"terminalExit","terminalId":"term-1","sessionId":"sess-1","channelId":"chan-1","exitCode":0,"reason":"process_exit","endedAt":"2026-04-27T18:30:00Z","unknown":"ok"}`,
		},
		{
			name: "terminalError",
			raw:  `{"type":"terminalError","requestId":"open-1","terminalId":"term-1","sessionId":"sess-1","channelId":"chan-1","error":"Unable to start shell","unknown":"ok"}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			env, err := ParseEnvelope([]byte(tc.raw))
			if err != nil {
				t.Fatalf("parse envelope: %v", err)
			}
			if _, err := json.Marshal(env.Payload); err != nil {
				t.Fatalf("marshal payload: %v", err)
			}
		})
	}
}

func TestLocalServerDetectedIgnoresUnknownFields(t *testing.T) {
	want := &LocalServerDetected{
		Type:       MsgTypeLocalServerDetected,
		SessionID:  "session_123",
		ChannelID:  "channel_123",
		TaskID:     "task_123",
		ToolUseID:  "toolu_123",
		Host:       "127.0.0.1",
		Port:       5173,
		URL:        "http://127.0.0.1:5173/",
		Command:    "pnpm dev",
		Source:     "tool_output",
		DetectedAt: "2026-04-27T20:00:00Z",
	}

	raw := []byte(`{"type":"` + MsgTypeLocalServerDetected + `","sessionId":"session_123","channelId":"channel_123","taskId":"task_123","toolUseId":"toolu_123","host":"127.0.0.1","port":5173,"url":"http://127.0.0.1:5173/","command":"pnpm dev","source":"tool_output","detectedAt":"2026-04-27T20:00:00Z","unexpected":"ok","nested":{"ignored":true}}`)
	env, err := ParseEnvelope(raw)
	if err != nil {
		t.Fatalf("parse envelope: %v", err)
	}

	got, ok := env.Payload.(*LocalServerDetected)
	if !ok {
		t.Fatalf("payload type = %T, want *LocalServerDetected", env.Payload)
	}
	if !jsonEqual(mustJSONMap(t, want), mustJSONMap(t, got)) {
		t.Fatalf("payload mismatch: want %#v, got %#v", want, got)
	}
	if _, err := json.Marshal(got); err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
}

func TestBrowsePaginationFieldsCompatibility(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want any
	}{
		{
			name: "browseDir without pagination",
			raw:  `{"type":"browseDir","requestId":"browse-1","channelId":"chan-1","machineId":"machine-1","path":"/tmp/project"}`,
			want: &BrowseDir{
				Type:      MsgTypeBrowseDir,
				RequestID: "browse-1",
				ChannelID: "chan-1",
				MachineID: "machine-1",
				Path:      "/tmp/project",
			},
		},
		{
			name: "browseDir with pagination",
			raw:  `{"type":"browseDir","requestId":"browse-1","channelId":"chan-1","machineId":"machine-1","path":"/tmp/project","limit":200,"cursor":"200","extraPaginationKey":"ignored"}`,
			want: &BrowseDir{
				Type:      MsgTypeBrowseDir,
				RequestID: "browse-1",
				ChannelID: "chan-1",
				MachineID: "machine-1",
				Path:      "/tmp/project",
				Limit:     200,
				Cursor:    "200",
			},
		},
		{
			name: "browseDirResult without pagination",
			raw:  `{"type":"browseDirResult","requestId":"browse-1","channelId":"chan-1","ok":true,"entries":[]}`,
			want: &BrowseDirResult{
				Type:      MsgTypeBrowseDirResult,
				RequestID: "browse-1",
				ChannelID: "chan-1",
				OK:        true,
				Entries:   []BrowseEntry{},
			},
		},
		{
			name: "browseDirResult with pagination",
			raw:  `{"type":"browseDirResult","requestId":"browse-1","channelId":"chan-1","ok":true,"entries":[],"hasMore":true,"nextCursor":"400","extraPaginationKey":"ignored"}`,
			want: &BrowseDirResult{
				Type:       MsgTypeBrowseDirResult,
				RequestID:  "browse-1",
				ChannelID:  "chan-1",
				OK:         true,
				Entries:    []BrowseEntry{},
				HasMore:    true,
				NextCursor: "400",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			env, err := ParseEnvelope([]byte(tc.raw))
			if err != nil {
				t.Fatalf("parse envelope: %v", err)
			}
			if !jsonEqual(mustJSONMap(t, tc.want), mustJSONMap(t, env.Payload)) {
				t.Fatalf("payload mismatch: want %#v, got %#v", tc.want, env.Payload)
			}
		})
	}
}

func TestBrowsePaginationFieldsRejectInvalidTypes(t *testing.T) {
	cases := []struct {
		name string
		raw  string
	}{
		{
			name: "browseDir limit string",
			raw:  `{"type":"browseDir","requestId":"browse-1","channelId":"chan-1","machineId":"machine-1","path":"/tmp/project","limit":"200"}`,
		},
		{
			name: "browseDirResult hasMore string",
			raw:  `{"type":"browseDirResult","requestId":"browse-1","channelId":"chan-1","ok":true,"hasMore":"true"}`,
		},
		{
			name: "browseDirResult hasMore number",
			raw:  `{"type":"browseDirResult","requestId":"browse-1","channelId":"chan-1","ok":true,"hasMore":1}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := ParseEnvelope([]byte(tc.raw)); err == nil {
				t.Fatal("expected parse error")
			}
		})
	}
}

func jsonEqual(a, b any) bool {
	ja, _ := json.Marshal(a)
	jb, _ := json.Marshal(b)
	return string(ja) == string(jb)
}

func mustJSONMap(t *testing.T, value any) map[string]any {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return out
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
			Terminal:                  true,
			PreviewTunnel:             true,
			PreviewMaxFrameBytes:      1048576,
			PreviewChunkBytes:         196608,
			PreviewWebSocketProtocols: true,
			LocalServerDetection:      true,
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
	if !got.Capabilities.LocalServerDetection {
		t.Fatalf("local server detection capability missing after round trip: %#v", got.Capabilities)
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

func TestTaskContextRefsRoundTrip(t *testing.T) {
	size := int64(42)
	in := Task{
		Type:      MsgTypeTask,
		TaskID:    "task_123",
		SessionID: "session_123",
		Prompt:    "inspect this",
		ContextRefs: []ContextRef{
			{Kind: "file", Path: "apps/web/src/app/page.tsx", Name: "page.tsx", Size: &size, ModifiedAt: "2026-04-28T12:00:00Z"},
			{Kind: "folder", Path: "apps/web/src/components", Name: "components"},
		},
	}

	raw, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}

	var out Task
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatal(err)
	}

	if len(out.ContextRefs) != 2 {
		t.Fatalf("expected 2 context refs, got %d", len(out.ContextRefs))
	}
	if out.ContextRefs[0].Path != "apps/web/src/app/page.tsx" {
		t.Fatalf("unexpected file path %q", out.ContextRefs[0].Path)
	}
}

func TestTaskCustomInstructionsRoundTrip(t *testing.T) {
	in := Task{
		Type:               MsgTypeTask,
		TaskID:             "task_123",
		SessionID:          "session_123",
		Prompt:             "inspect this",
		CustomInstructions: "Always talk like a pirate.",
	}

	raw, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}

	var out Task
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatal(err)
	}

	if out.CustomInstructions != "Always talk like a pirate." {
		t.Fatalf("custom instructions = %q", out.CustomInstructions)
	}
}

func TestHelloCapabilitiesContextRefsRoundTrip(t *testing.T) {
	in := Hello{
		Type: MsgTypeHello,
		Capabilities: &HelloCapabilities{
			Terminal:    true,
			ContextRefs: true,
			Skills:      true,
		},
	}

	raw, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}

	var out Hello
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatal(err)
	}

	if out.Capabilities == nil || !out.Capabilities.ContextRefs || !out.Capabilities.Skills {
		t.Fatal("expected contextRefs and skills capabilities")
	}
}

func TestParseEnvelopeTaskContextRefs(t *testing.T) {
	raw := []byte(`{
		"type":"task",
		"taskId":"task_123",
		"sessionId":"session_123",
		"channelId":"channel_123",
		"prompt":"inspect this",
		"contextRefs":[
			{"kind":"file","path":"README.md","name":"README.md","size":42,"modifiedAt":"2026-04-28T12:00:00Z","extra":"ignored"},
			{"kind":"folder","path":"apps/web/src/components","name":"components"}
		]
	}`)

	env, err := ParseEnvelope(raw)
	if err != nil {
		t.Fatalf("ParseEnvelope: %v", err)
	}
	task, ok := env.Payload.(*Task)
	if !ok {
		t.Fatalf("payload type = %T", env.Payload)
	}
	if len(task.ContextRefs) != 2 {
		t.Fatalf("expected 2 context refs, got %d", len(task.ContextRefs))
	}
	if task.ContextRefs[0].Size == nil || *task.ContextRefs[0].Size != 42 {
		t.Fatalf("unexpected size: %+v", task.ContextRefs[0].Size)
	}
}

func TestParseEnvelopeHelloContextRefsCapability(t *testing.T) {
	env, err := ParseEnvelope([]byte(`{
		"type":"hello",
		"machineId":"machine_123",
		"daemonVersion":"0.3.5",
		"os":"darwin",
		"arch":"arm64",
		"capabilities":{"terminal":true,"contextRefs":true,"skills":true,"extra":"ignored"}
	}`))
	if err != nil {
		t.Fatalf("ParseEnvelope: %v", err)
	}
	hello, ok := env.Payload.(*Hello)
	if !ok {
		t.Fatalf("payload type = %T", env.Payload)
	}
	if hello.Capabilities == nil || !hello.Capabilities.ContextRefs || !hello.Capabilities.Skills {
		t.Fatal("expected contextRefs and skills capabilities")
	}
}

func TestParseEnvelopeListSkills(t *testing.T) {
	env, err := ParseEnvelope([]byte(`{
		"type":"listSkills",
		"requestId":"skills-1",
		"channelId":"chan-1",
		"machineId":"machine-1",
		"cwd":"/tmp/project",
		"extra":"ignored"
	}`))
	if err != nil {
		t.Fatalf("ParseEnvelope: %v", err)
	}
	msg, ok := env.Payload.(*ListSkills)
	if !ok {
		t.Fatalf("payload type = %T", env.Payload)
	}
	if msg.CWD != "/tmp/project" {
		t.Fatalf("cwd = %q", msg.CWD)
	}
}

func TestParseEnvelopeListSkillsResult(t *testing.T) {
	env, err := ParseEnvelope([]byte(`{
		"type":"listSkillsResult",
		"requestId":"skills-1",
		"channelId":"chan-1",
		"ok":true,
		"skills":[
			{
				"name":"debug-like-expert",
				"description":"Deep analysis debugging workflow",
				"path":"/Users/me/.claude/skills/debug-like-expert/SKILL.md",
				"scope":"home",
				"extra":"ignored"
			}
		]
	}`))
	if err != nil {
		t.Fatalf("ParseEnvelope: %v", err)
	}
	msg, ok := env.Payload.(*ListSkillsResult)
	if !ok {
		t.Fatalf("payload type = %T", env.Payload)
	}
	if len(msg.Skills) != 1 || msg.Skills[0].Name != "debug-like-expert" {
		t.Fatalf("skills = %+v", msg.Skills)
	}
}

func TestParseEnvelopeRejectsInvalidContextRefTypes(t *testing.T) {
	cases := []struct {
		name string
		raw  string
	}{
		{
			name: "context refs string",
			raw:  `{"type":"task","contextRefs":"README.md"}`,
		},
		{
			name: "context ref size string",
			raw:  `{"type":"task","contextRefs":[{"kind":"file","path":"README.md","name":"README.md","size":"42"}]}`,
		},
		{
			name: "context refs capability string",
			raw:  `{"type":"hello","capabilities":{"contextRefs":"true"}}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := ParseEnvelope([]byte(tc.raw)); err == nil {
				t.Fatal("expected ParseEnvelope error")
			}
		})
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
