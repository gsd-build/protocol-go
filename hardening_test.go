package protocol

import (
	"strings"
	"testing"
)

func TestParseEnvelopeWithLimitsRejectsOversizedFrame(t *testing.T) {
	raw := []byte(`{"type":"task","taskId":"task_1","sessionId":"session_1","channelId":"channel_1","prompt":"` + strings.Repeat("x", 64) + `"}`)

	_, err := ParseEnvelopeWithLimits(raw, EnvelopeLimits{MaxFrameBytes: 64})
	if err == nil {
		t.Fatal("expected oversized frame error")
	}
}

func TestParseEnvelopeWithLimitsRejectsExcessiveDepth(t *testing.T) {
	raw := []byte(`{"type":"stream","sessionId":"session_1","channelId":"channel_1","event":{"a":{"b":{"c":true}}}}`)

	_, err := ParseEnvelopeWithLimits(raw, EnvelopeLimits{MaxDepth: 3})
	if err == nil {
		t.Fatal("expected depth limit error")
	}
}

func TestParseEnvelopeWithLimitsRejectsTooManyObjectFields(t *testing.T) {
	raw := []byte(`{"type":"hello","machineId":"machine_1","daemonVersion":"1.0.0","os":"darwin","arch":"arm64"}`)

	_, err := ParseEnvelopeWithLimits(raw, EnvelopeLimits{MaxObjectFields: 2})
	if err == nil {
		t.Fatal("expected object field limit error")
	}
}

func TestParseEnvelopeWithLimitsRejectsTooManyArrayItems(t *testing.T) {
	raw := []byte(`{"type":"task","taskId":"task_1","sessionId":"session_1","channelId":"channel_1","imageUrls":["a","b","c"]}`)

	_, err := ParseEnvelopeWithLimits(raw, EnvelopeLimits{MaxArrayItems: 2})
	if err == nil {
		t.Fatal("expected array item limit error")
	}
}

func TestValidateEnvelopeFrameRejectsNonObjectRoot(t *testing.T) {
	err := ValidateEnvelopeFrame([]byte(`[{"type":"hello"}]`), EnvelopeLimits{})
	if err == nil {
		t.Fatal("expected root shape error")
	}
}

func TestParseEnvelopeWithLimitsAcceptsBoundedFrame(t *testing.T) {
	raw := []byte(`{"type":"taskStarted","taskId":"task_1","sessionId":"session_1","channelId":"channel_1","startedAt":"2026-04-28T12:00:00Z","requestId":"request_1"}`)

	env, err := ParseEnvelopeWithLimits(raw, EnvelopeLimits{
		MaxFrameBytes:   1024,
		MaxDepth:        4,
		MaxObjectFields: 8,
		MaxArrayItems:   4,
	})
	if err != nil {
		t.Fatalf("ParseEnvelopeWithLimits: %v", err)
	}
	if env.Type != MsgTypeTaskStarted {
		t.Fatalf("type = %q, want %q", env.Type, MsgTypeTaskStarted)
	}
}

func TestValidateRequestBinding(t *testing.T) {
	request := &Task{
		Type:      MsgTypeTask,
		TaskID:    "task_1",
		SessionID: "session_1",
		ChannelID: "channel_1",
		RequestID: "request_1",
	}
	response := &TaskStarted{
		Type:      MsgTypeTaskStarted,
		TaskID:    "task_1",
		SessionID: "session_1",
		ChannelID: "channel_1",
		RequestID: "request_1",
	}
	if err := ValidateRequestBinding(request, response); err != nil {
		t.Fatalf("ValidateRequestBinding: %v", err)
	}

	response.RequestID = "request_2"
	if err := ValidateRequestBinding(request, response); err == nil {
		t.Fatal("expected requestId mismatch")
	}

	response.RequestID = ""
	if err := ValidateRequestBinding(request, response); err == nil {
		t.Fatal("expected missing actual requestId")
	}
}

func TestValidateSessionBinding(t *testing.T) {
	expected := Binding{SessionID: "session_1", ChannelID: "channel_1"}
	actual := &TaskError{
		Type:      MsgTypeTaskError,
		TaskID:    "task_1",
		SessionID: "session_1",
		ChannelID: "channel_1",
		Error:     "failed",
	}
	if err := ValidateSessionBinding(expected, actual); err != nil {
		t.Fatalf("ValidateSessionBinding: %v", err)
	}

	actual.ChannelID = "channel_2"
	if err := ValidateSessionBinding(expected, actual); err == nil {
		t.Fatal("expected channelId mismatch")
	}

	actual.ChannelID = "channel_1"
	actual.SessionID = "session_2"
	if err := ValidateSessionBinding(expected, actual); err == nil {
		t.Fatal("expected sessionId mismatch")
	}
}

func TestExtractBindingFromEnvelope(t *testing.T) {
	env := &Envelope{
		Type: MsgTypePreviewHTTPRequest,
		Payload: &PreviewHTTPRequest{
			Type:      MsgTypePreviewHTTPRequest,
			RequestID: "request_1",
			StreamID:  "stream_1",
			PreviewID: "preview_1",
		},
	}

	binding := ExtractBinding(env)
	if binding.RequestID != "request_1" || binding.StreamID != "stream_1" || binding.PreviewID != "preview_1" {
		t.Fatalf("unexpected binding: %+v", binding)
	}
}
