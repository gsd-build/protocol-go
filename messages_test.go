package protocol

import (
	"encoding/json"
	"testing"
)

func TestEnvelopeRoundTrip(t *testing.T) {
	cases := []struct {
		name string
		msg  any
	}{
		{"task", &Task{
			Type:           MsgTypeTask,
			TaskID:         "11111111-1111-1111-1111-111111111111",
			SessionID:      "22222222-2222-2222-2222-222222222222",
			ChannelID:      "ch-1",
			Prompt:         "hello",
			Model:          "claude-opus-4-6[1m]",
			Effort:         "max",
			PermissionMode: "acceptEdits",
			CWD:            "/tmp/project",
		}},
		{"stream", &Stream{
			Type:           MsgTypeStream,
			SessionID:      "22222222-2222-2222-2222-222222222222",
			ChannelID:      "ch-1",
			SequenceNumber: 42,
			Event:          json.RawMessage(`{"delta":{"text":"hi"}}`),
		}},
		{"ack", &Ack{
			Type:           MsgTypeAck,
			SessionID:      "33333333-3333-3333-3333-333333333333",
			SequenceNumber: 99,
		}},
		{"hello", &Hello{
			Type:          MsgTypeHello,
			MachineID:     "mach-id",
			DaemonVersion: "0.1.0",
			OS:            "darwin",
			Arch:          "arm64",
			LastSequenceBySession: map[string]int64{
				"22222222-2222-2222-2222-222222222222": 42,
			},
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
