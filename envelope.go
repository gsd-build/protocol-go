package protocol

import (
	"encoding/json"
	"fmt"
)

// Envelope is a parsed message ready for type-switching.
type Envelope struct {
	Type    string
	Payload any
}

// ParseEnvelope reads raw JSON, looks at the type field, and unmarshals
// into the correct concrete struct.
func ParseEnvelope(data []byte) (*Envelope, error) {
	var peek struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return nil, fmt.Errorf("envelope: %w", err)
	}

	var payload any
	switch peek.Type {
	case MsgTypeTask:
		payload = &Task{}
	case MsgTypeStop:
		payload = &Stop{}
	case MsgTypePermissionResponse:
		payload = &PermissionResponse{}
	case MsgTypeQuestionResponse:
		payload = &QuestionResponse{}
	case MsgTypeBrowseDir:
		payload = &BrowseDir{}
	case MsgTypeReadFile:
		payload = &ReadFile{}
	case MsgTypeStream:
		payload = &Stream{}
	case MsgTypeTaskStarted:
		payload = &TaskStarted{}
	case MsgTypeTaskComplete:
		payload = &TaskComplete{}
	case MsgTypeTaskError:
		payload = &TaskError{}
	case MsgTypePermissionRequest:
		payload = &PermissionRequest{}
	case MsgTypeQuestion:
		payload = &Question{}
	case MsgTypeHeartbeat:
		payload = &Heartbeat{}
	case MsgTypeBrowseDirResult:
		payload = &BrowseDirResult{}
	case MsgTypeReadFileResult:
		payload = &ReadFileResult{}
	case MsgTypeHello:
		payload = &Hello{}
	case MsgTypeWelcome:
		payload = &Welcome{}
	case MsgTypeAck:
		payload = &Ack{}
	case MsgTypeReplayRequest:
		payload = &ReplayRequest{}
	default:
		return nil, fmt.Errorf("unknown message type: %q", peek.Type)
	}

	if err := json.Unmarshal(data, payload); err != nil {
		return nil, fmt.Errorf("unmarshal %s: %w", peek.Type, err)
	}

	return &Envelope{
		Type:    peek.Type,
		Payload: payload,
	}, nil
}
