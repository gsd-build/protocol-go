package protocol

import (
	"encoding/json"
	"reflect"
	"testing"
)

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
		{"syncCrons", &SyncCrons{
			Type:      MsgTypeSyncCrons,
			MachineID: "mach-id",
			SentAt:    "2026-04-14T12:00:00Z",
			Jobs: []CronSpec{
				{
					ID:              "11111111-1111-1111-1111-111111111111",
					Name:            "Daily tests",
					CronExpression:  "0 9 * * *",
					Prompt:          "Run tests",
					Mode:            "fresh",
					Model:           "claude-sonnet-4-6",
					Effort:          "medium",
					ProjectID:       "22222222-2222-2222-2222-222222222222",
					TargetSessionID: "33333333-3333-3333-3333-333333333333",
					Enabled:         true,
				},
			},
		}},
		{"cronInventory", &CronInventory{
			Type:      MsgTypeCronInventory,
			MachineID: "mach-id",
			Timestamp: "2026-04-14T12:00:00Z",
			Items: []CronLocalState{
				{
					ID:              "11111111-1111-1111-1111-111111111111",
					Name:            "Daily tests",
					CronExpression:  "0 9 * * *",
					Enabled:         true,
					SyncedAt:        "2026-04-14T12:00:00Z",
					LastRunAt:       "2026-04-14T12:05:00Z",
					NextRunAt:       "2026-04-15T09:00:00Z",
					LocallyModified: false,
				},
			},
		}},
		{"cronExecResult", &CronExecResult{
			Type:       MsgTypeCronExecResult,
			MachineID:  "mach-id",
			CronJobID:  "11111111-1111-1111-1111-111111111111",
			TaskID:     "44444444-4444-4444-4444-444444444444",
			Status:     "success",
			DurationMs: 45000,
			Timestamp:  "2026-04-14T12:00:00Z",
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
		{"skillInventory", &SkillInventory{
			Type:      MsgTypeSkillInventory,
			MachineID: "mach-id",
			Entries: []SkillInventoryEntry{
				{
					Slug:               "write-plan",
					DisplayName:        "Write Plan",
					Description:        "Generate implementation plans",
					Scope:              "global",
					Runtime:            "claude",
					Root:               "/Users/lex/.claude/skills",
					ProjectRoot:        "/Users/lex/Developer/gsd",
					RelativePath:       "write-plan/SKILL.md",
					SourceKind:         "skill_dir",
					MachineFingerprint: "machine:abc123",
					Editable:           true,
				},
			},
		}},
		{"skillContentRequest", &SkillContentRequest{
			Type:         MsgTypeSkillContentRequest,
			MachineID:    "mach-id",
			Slug:         "write-plan",
			Root:         "/Users/lex/.claude/skills",
			RelativePath: "write-plan/SKILL.md",
		}},
		{"skillContentUpload", &SkillContentUpload{
			Type:               MsgTypeSkillContentUpload,
			MachineID:          "mach-id",
			Slug:               "write-plan",
			Root:               "/Users/lex/.claude/skills",
			RelativePath:       "write-plan/SKILL.md",
			Content:            "---\nname: Write Plan\n---\n",
			MachineFingerprint: "machine:abc123",
			BaseCloudRevision:  7,
		}},
		{"skillPush", &SkillPush{
			Type:             MsgTypeSkillPush,
			MachineID:        "mach-id",
			Slug:             "write-plan",
			Root:             "/Users/lex/.claude/skills",
			RelativePath:     "write-plan/SKILL.md",
			Content:          "---\nname: Write Plan\n---\n",
			CloudFingerprint: "cloud:def456",
			CloudRevision:    8,
		}},
		{"skillDelete", &SkillDelete{
			Type:          MsgTypeSkillDelete,
			MachineID:     "mach-id",
			Slug:          "write-plan",
			Root:          "/Users/lex/.claude/skills",
			RelativePath:  "write-plan/SKILL.md",
			CloudRevision: 8,
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
			_ = json.Unmarshal(data, &original)
			_ = json.Unmarshal(reMarshaled, &final)

			for k, v := range original {
				got, ok := final[k]
				if !ok {
					t.Errorf("missing key %q after round trip", k)
					continue
				}
				if !jsonEqual(v, got) {
					t.Errorf("value mismatch for %q: want %v, got %v", k, v, got)
				}
			}
		})
	}
}

func TestSkillMessageConcretePayloadTypes(t *testing.T) {
	cases := []struct {
		name       string
		msg        any
		wantType   any
		wantFields func(t *testing.T, payload any)
	}{
		{
			name: "skillInventory",
			msg: &SkillInventory{
				Type:      MsgTypeSkillInventory,
				MachineID: "mach-id",
				Entries: []SkillInventoryEntry{
					{
						Slug:               "write-plan",
						DisplayName:        "Write Plan",
						Description:        "Generate implementation plans",
						Scope:              "global",
						Runtime:            "claude",
						Root:               "/Users/lex/.claude/skills",
						ProjectRoot:        "/Users/lex/Developer/gsd",
						RelativePath:       "write-plan/SKILL.md",
						SourceKind:         "skill_dir",
						MachineFingerprint: "machine:abc123",
						Editable:           true,
					},
				},
			},
			wantType: (*SkillInventory)(nil),
			wantFields: func(t *testing.T, payload any) {
				t.Helper()
				got := payload.(*SkillInventory)
				if got.Entries[0].ProjectRoot != "/Users/lex/Developer/gsd" {
					t.Fatalf("expected projectRoot to survive round-trip, got %q", got.Entries[0].ProjectRoot)
				}
				if got.Entries[0].MachineFingerprint != "machine:abc123" {
					t.Fatalf("expected machine fingerprint to survive round-trip, got %q", got.Entries[0].MachineFingerprint)
				}
			},
		},
		{
			name: "skillContentRequest",
			msg: &SkillContentRequest{
				Type:         MsgTypeSkillContentRequest,
				MachineID:    "mach-id",
				Slug:         "write-plan",
				Root:         "/Users/lex/.claude/skills",
				RelativePath: "write-plan/SKILL.md",
			},
			wantType: (*SkillContentRequest)(nil),
		},
		{
			name: "skillContentUpload",
			msg: &SkillContentUpload{
				Type:               MsgTypeSkillContentUpload,
				MachineID:          "mach-id",
				Slug:               "write-plan",
				Root:               "/Users/lex/.claude/skills",
				RelativePath:       "write-plan/SKILL.md",
				Content:            "---\nname: Write Plan\n---\n",
				MachineFingerprint: "machine:abc123",
				BaseCloudRevision:  7,
			},
			wantType: (*SkillContentUpload)(nil),
			wantFields: func(t *testing.T, payload any) {
				t.Helper()
				got := payload.(*SkillContentUpload)
				if got.BaseCloudRevision != 7 {
					t.Fatalf("expected base cloud revision to survive round-trip, got %d", got.BaseCloudRevision)
				}
				if got.MachineFingerprint != "machine:abc123" {
					t.Fatalf("expected machine fingerprint to survive round-trip, got %q", got.MachineFingerprint)
				}
			},
		},
		{
			name: "skillPush",
			msg: &SkillPush{
				Type:             MsgTypeSkillPush,
				MachineID:        "mach-id",
				Slug:             "write-plan",
				Root:             "/Users/lex/.claude/skills",
				RelativePath:     "write-plan/SKILL.md",
				Content:          "---\nname: Write Plan\n---\n",
				CloudFingerprint: "cloud:def456",
				CloudRevision:    8,
			},
			wantType: (*SkillPush)(nil),
			wantFields: func(t *testing.T, payload any) {
				t.Helper()
				got := payload.(*SkillPush)
				if got.CloudRevision != 8 {
					t.Fatalf("expected cloud revision to survive round-trip, got %d", got.CloudRevision)
				}
				if got.CloudFingerprint != "cloud:def456" {
					t.Fatalf("expected cloud fingerprint to survive round-trip, got %q", got.CloudFingerprint)
				}
			},
		},
		{
			name: "skillDelete",
			msg: &SkillDelete{
				Type:          MsgTypeSkillDelete,
				MachineID:     "mach-id",
				Slug:          "write-plan",
				Root:          "/Users/lex/.claude/skills",
				RelativePath:  "write-plan/SKILL.md",
				CloudRevision: 8,
			},
			wantType: (*SkillDelete)(nil),
			wantFields: func(t *testing.T, payload any) {
				t.Helper()
				got := payload.(*SkillDelete)
				if got.CloudRevision != 8 {
					t.Fatalf("expected cloud revision to survive round-trip, got %d", got.CloudRevision)
				}
			},
		},
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

			if gotType := reflect.TypeOf(env.Payload); gotType != reflect.TypeOf(tc.wantType) {
				t.Fatalf("unexpected payload type: got %v, want %v", gotType, reflect.TypeOf(tc.wantType))
			}

			if tc.wantFields != nil {
				tc.wantFields(t, env.Payload)
			}
		})
	}
}

func TestSkillInventoryEmptyWireShape(t *testing.T) {
	msg := &SkillInventory{
		Type:      MsgTypeSkillInventory,
		MachineID: "mach-id",
		Entries:   []SkillInventoryEntry{},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal raw: %v", err)
	}

	entries, ok := raw["entries"].([]any)
	if !ok {
		t.Fatalf("expected entries to be an array, got %T", raw["entries"])
	}
	if len(entries) != 0 {
		t.Fatalf("expected empty entries array, got %d items", len(entries))
	}

	env, err := ParseEnvelope(data)
	if err != nil {
		t.Fatalf("parse envelope: %v", err)
	}

	got, ok := env.Payload.(*SkillInventory)
	if !ok {
		t.Fatalf("expected *SkillInventory, got %T", env.Payload)
	}
	if len(got.Entries) != 0 {
		t.Fatalf("expected zero entries after parse, got %d", len(got.Entries))
	}
}

func jsonEqual(a, b any) bool {
	ja, _ := json.Marshal(a)
	jb, _ := json.Marshal(b)
	return string(ja) == string(jb)
}

func TestParseEnvelopeRejectsUnknownType(t *testing.T) {
	_, err := ParseEnvelope([]byte(`{"type":"bogus"}`))
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestHelloRoundTripWithActiveTasks(t *testing.T) {
	h := &Hello{
		Type:          MsgTypeHello,
		MachineID:     "m-1",
		DaemonVersion: "0.2.0",
		OS:            "darwin",
		Arch:          "arm64",
		ActiveTasks:   []string{"task-a", "task-b"},
	}

	data, err := json.Marshal(h)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	env, err := ParseEnvelope(data)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	got, ok := env.Payload.(*Hello)
	if !ok {
		t.Fatalf("expected *Hello, got %T", env.Payload)
	}
	if len(got.ActiveTasks) != 2 {
		t.Fatalf("expected 2 active tasks, got %d", len(got.ActiveTasks))
	}
	if got.ActiveTasks[0] != "task-a" || got.ActiveTasks[1] != "task-b" {
		t.Errorf("unexpected active tasks: %v", got.ActiveTasks)
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
	_ = json.Unmarshal(data, &raw)
	if _, exists := raw["activeTasks"]; exists {
		t.Error("activeTasks should be omitted when empty")
	}
}

func TestWelcomeRoundTripWithVersion(t *testing.T) {
	w := &Welcome{
		Type:                MsgTypeWelcome,
		LatestDaemonVersion: "0.2.1",
	}

	data, err := json.Marshal(w)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	env, err := ParseEnvelope(data)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	got, ok := env.Payload.(*Welcome)
	if !ok {
		t.Fatalf("expected *Welcome, got %T", env.Payload)
	}
	if got.LatestDaemonVersion != "0.2.1" {
		t.Errorf("expected 0.2.1, got %s", got.LatestDaemonVersion)
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
	_ = json.Unmarshal(data, &raw)
	if _, exists := raw["latestDaemonVersion"]; exists {
		t.Error("latestDaemonVersion should be omitted when empty")
	}
}

func TestUpdateAvailableRoundTrip(t *testing.T) {
	ua := &UpdateAvailable{
		Type:           MsgTypeUpdateAvailable,
		CurrentVersion: "0.1.0",
		LatestVersion:  "0.2.0",
	}

	data, err := json.Marshal(ua)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	env, err := ParseEnvelope(data)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	got, ok := env.Payload.(*UpdateAvailable)
	if !ok {
		t.Fatalf("expected *UpdateAvailable, got %T", env.Payload)
	}
	if got.CurrentVersion != "0.1.0" {
		t.Errorf("expected 0.1.0, got %s", got.CurrentVersion)
	}
	if got.LatestVersion != "0.2.0" {
		t.Errorf("expected 0.2.0, got %s", got.LatestVersion)
	}
}

func TestRequestIDRoundTrip(t *testing.T) {
	requestID := "33333333-3333-3333-3333-333333333333"
	tp := "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"

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
		RequestID:      requestID,
		Traceparent:    tp,
	}

	data, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	env, err := ParseEnvelope(data)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	got, ok := env.Payload.(*Task)
	if !ok {
		t.Fatalf("expected *Task, got %T", env.Payload)
	}
	if got.RequestID != requestID {
		t.Errorf("requestId mismatch: want %s, got %s", requestID, got.RequestID)
	}
	if got.Traceparent != tp {
		t.Errorf("traceparent mismatch: want %s, got %s", tp, got.Traceparent)
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
	_ = json.Unmarshal(data, &raw)
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
	ms := &MachineStatus{
		Type:      MsgTypeMachineStatus,
		MachineID: "m-1",
		Online:    true,
	}

	data, err := json.Marshal(ms)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	env, err := ParseEnvelope(data)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	got, ok := env.Payload.(*MachineStatus)
	if !ok {
		t.Fatalf("expected *MachineStatus, got %T", env.Payload)
	}
	if got.MachineID != "m-1" {
		t.Errorf("expected m-1, got %s", got.MachineID)
	}
	if !got.Online {
		t.Error("expected online=true")
	}
}
