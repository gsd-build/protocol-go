package protocol

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// Envelope is a parsed message ready for type-switching.
type Envelope struct {
	Type    string
	Payload any
}

func (e Envelope) DecodePayload() (any, error) {
	payload, err := payloadForType(e.Type)
	if err != nil {
		return nil, err
	}
	switch raw := e.Payload.(type) {
	case json.RawMessage:
		if err := json.Unmarshal(raw, payload); err != nil {
			return nil, fmt.Errorf("decode %s: %w", e.Type, err)
		}
	case []byte:
		if err := json.Unmarshal(raw, payload); err != nil {
			return nil, fmt.Errorf("decode %s: %w", e.Type, err)
		}
	default:
		return e.Payload, nil
	}
	return payload, nil
}

// ParseEnvelope reads raw JSON, looks at the type field, and unmarshals
// into the correct concrete struct.
func ParseEnvelope(data []byte) (*Envelope, error) {
	return parseEnvelope(data)
}

// ParseEnvelopeWithLimits validates frame size, JSON nesting depth, object
// field counts, and array element counts before parsing the envelope.
func ParseEnvelopeWithLimits(data []byte, limits EnvelopeLimits) (*Envelope, error) {
	if err := ValidateEnvelopeFrame(data, limits); err != nil {
		return nil, err
	}
	return parseEnvelope(data)
}

func parseEnvelope(data []byte) (*Envelope, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("envelope: %w", err)
	}

	rawType, ok := raw["type"]
	if !ok {
		return nil, fmt.Errorf("envelope: missing type")
	}

	var msgType string
	if err := json.Unmarshal(rawType, &msgType); err != nil {
		return nil, fmt.Errorf("envelope type: %w", err)
	}

	payload, err := payloadForType(msgType)
	if err != nil {
		return nil, err
	}

	if err := populatePayload(payload, raw); err != nil {
		return nil, fmt.Errorf("unmarshal %s: %w", msgType, err)
	}

	return &Envelope{
		Type:    msgType,
		Payload: payload,
	}, nil
}

func payloadForType(msgType string) (any, error) {
	switch msgType {
	case MsgTypeTask:
		return &Task{}, nil
	case MsgTypeTaskLifecycle:
		return &TaskLifecycle{}, nil
	case MsgTypeStop:
		return &Stop{}, nil
	case MsgTypePermissionResponse:
		return &PermissionResponse{}, nil
	case MsgTypeQuestionResponse:
		return &QuestionResponse{}, nil
	case MsgTypeBrowseDir:
		return &BrowseDir{}, nil
	case MsgTypeReadFile:
		return &ReadFile{}, nil
	case MsgTypeMkDir:
		return &MkDir{}, nil
	case MsgTypeMkDirResult:
		return &MkDirResult{}, nil
	case MsgTypeListSkills:
		return &ListSkills{}, nil
	case MsgTypeListSkillsResult:
		return &ListSkillsResult{}, nil
	case MsgTypeCompactRequest:
		return &CompactRequest{}, nil
	case MsgTypeContextStatsRequest:
		return &ContextStatsRequest{}, nil
	case MsgTypeStream:
		return &Stream{}, nil
	case MsgTypeTaskStarted:
		return &TaskStarted{}, nil
	case MsgTypeTaskComplete:
		return &TaskComplete{}, nil
	case MsgTypeTaskError:
		return &TaskError{}, nil
	case MsgTypeTaskCancelled:
		return &TaskCancelled{}, nil
	case MsgTypePermissionRequest:
		return &PermissionRequest{}, nil
	case MsgTypeQuestion:
		return &Question{}, nil
	case MsgTypeHeartbeat:
		return &Heartbeat{}, nil
	case MsgTypeBrowseDirResult:
		return &BrowseDirResult{}, nil
	case MsgTypeReadFileResult:
		return &ReadFileResult{}, nil
	case MsgTypeContextStats:
		return &ContextStats{}, nil
	case MsgTypeCompactStatus:
		return &CompactStatus{}, nil
	case MsgTypeHello:
		return &Hello{}, nil
	case MsgTypeWelcome:
		return &Welcome{}, nil
	case MsgTypeMachineStatus:
		return &MachineStatus{}, nil
	case MsgTypePreviewOpen:
		return &PreviewOpen{}, nil
	case MsgTypePreviewOpenResult:
		return &PreviewOpenResult{}, nil
	case MsgTypePreviewClose:
		return &PreviewClose{}, nil
	case MsgTypePreviewHTTPRequest:
		return &PreviewHTTPRequest{}, nil
	case MsgTypePreviewHTTPResponseHead:
		return &PreviewHTTPResponseHead{}, nil
	case MsgTypePreviewStreamChunk:
		return &PreviewStreamChunk{}, nil
	case MsgTypePreviewStreamCancel:
		return &PreviewStreamCancel{}, nil
	case MsgTypePreviewWebSocketOpen:
		return &PreviewWebSocketOpen{}, nil
	case MsgTypePreviewWebSocketOpenResult:
		return &PreviewWebSocketOpenResult{}, nil
	case MsgTypePreviewWebSocketData:
		return &PreviewWebSocketData{}, nil
	case MsgTypePreviewWebSocketClose:
		return &PreviewWebSocketClose{}, nil
	case MsgTypeLocalServerDetected:
		return &LocalServerDetected{}, nil
	case MsgTypeTerminalOpen:
		return &TerminalOpen{}, nil
	case MsgTypeTerminalOpened:
		return &TerminalOpened{}, nil
	case MsgTypeTerminalInput:
		return &TerminalInput{}, nil
	case MsgTypeTerminalOutput:
		return &TerminalOutput{}, nil
	case MsgTypeTerminalSnapshot:
		return &TerminalSnapshot{}, nil
	case MsgTypeTerminalResize:
		return &TerminalResize{}, nil
	case MsgTypeTerminalClose:
		return &TerminalClose{}, nil
	case MsgTypeTerminalExit:
		return &TerminalExit{}, nil
	case MsgTypeTerminalError:
		return &TerminalError{}, nil
	case MsgTypeAgentTerminalStarted:
		return &AgentTerminalStarted{}, nil
	case MsgTypeAgentTerminalUpdated:
		return &AgentTerminalUpdated{}, nil
	case MsgTypeAgentTerminalAttach:
		return &AgentTerminalAttach{}, nil
	case MsgTypeAgentTerminalSnapshotRequest:
		return &AgentTerminalSnapshotRequest{}, nil
	case MsgTypeBrowserSessionOpen:
		return &BrowserSessionOpen{}, nil
	case MsgTypeBrowserSessionOpened:
		return &BrowserSessionOpened{}, nil
	case MsgTypeBrowserSessionClose:
		return &BrowserSessionClose{}, nil
	case MsgTypeBrowserSessionClosed:
		return &BrowserSessionClosed{}, nil
	case MsgTypeBrowserSessionError:
		return &BrowserSessionError{}, nil
	case MsgTypeBrowserFrame:
		return &BrowserFrame{}, nil
	case MsgTypeBrowserRefs:
		return &BrowserRefs{}, nil
	case MsgTypeBrowserCursor:
		return &BrowserCursor{}, nil
	case MsgTypeBrowserNavigation:
		return &BrowserNavigation{}, nil
	case MsgTypeBrowserAction:
		return &BrowserAction{}, nil
	case MsgTypeBrowserToolCall:
		return &BrowserToolCall{}, nil
	case MsgTypeBrowserToolResult:
		return &BrowserToolResult{}, nil
	case MsgTypeBrowserToolCallStarted:
		return &BrowserToolCallStarted{}, nil
	case MsgTypeBrowserToolCallUpdated:
		return &BrowserToolCallUpdated{}, nil
	case MsgTypeBrowserArtifactCreated:
		return &BrowserArtifactCreated{}, nil
	case MsgTypeBrowserControlClaim:
		return &BrowserControlClaim{}, nil
	case MsgTypeBrowserControlRelease:
		return &BrowserControlRelease{}, nil
	case MsgTypeBrowserUserInput:
		return &BrowserUserInput{}, nil
	case MsgTypeBrowserUserInputAck:
		return &BrowserUserInputAck{}, nil
	case MsgTypeBrowserTransportStatus:
		return &BrowserTransportStatus{}, nil
	case MsgTypeBrowserBridgeAccessOpen:
		return &BrowserBridgeAccessOpen{}, nil
	case MsgTypeBrowserBridgeAccessOpened:
		return &BrowserBridgeAccessOpened{}, nil
	case MsgTypeBrowserBridgeAccessClose:
		return &BrowserBridgeAccessClose{}, nil
	case MsgTypeBrowserSensitiveActionRequest:
		return &BrowserSensitiveActionRequest{}, nil
	case MsgTypeBrowserSensitiveActionResponse:
		return &BrowserSensitiveActionResponse{}, nil
	case MsgTypePlanningEvent:
		return &PlanningEvent{}, nil
	case MsgTypePlanningEventAck:
		return &PlanningEventAck{}, nil
	default:
		return nil, fmt.Errorf("unknown message type: %q", msgType)
	}
}

func populatePayload(dst any, raw map[string]json.RawMessage) error {
	value := reflect.ValueOf(dst)
	if value.Kind() != reflect.Pointer || value.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("payload must be a pointer to a struct, got %T", dst)
	}
	return populateStruct(value.Elem(), raw)
}

func populateStruct(dst reflect.Value, raw map[string]json.RawMessage) error {
	dstType := dst.Type()
	for i := 0; i < dstType.NumField(); i++ {
		fieldType := dstType.Field(i)
		if !fieldType.IsExported() {
			continue
		}

		fieldName, skip := jsonFieldName(fieldType)
		if skip {
			continue
		}

		fieldData, ok := raw[fieldName]
		if !ok {
			continue
		}

		if err := json.Unmarshal(fieldData, dst.Field(i).Addr().Interface()); err != nil {
			return fmt.Errorf("%s: %w", fieldName, err)
		}
	}
	return nil
}

func jsonFieldName(field reflect.StructField) (string, bool) {
	tag := field.Tag.Get("json")
	if tag == "-" {
		return "", true
	}

	if tag != "" {
		name, _, _ := strings.Cut(tag, ",")
		if name == "" {
			return field.Name, false
		}
		return name, false
	}

	return field.Name, false
}
