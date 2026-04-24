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

// ParseEnvelope reads raw JSON, looks at the type field, and unmarshals
// into the correct concrete struct.
func ParseEnvelope(data []byte) (*Envelope, error) {
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
	case MsgTypeHello:
		return &Hello{}, nil
	case MsgTypeWelcome:
		return &Welcome{}, nil
	case MsgTypeSyncCrons:
		return &SyncCrons{}, nil
	case MsgTypeCronInventory:
		return &CronInventory{}, nil
	case MsgTypeCronExecResult:
		return &CronExecResult{}, nil
	case MsgTypeMachineStatus:
		return &MachineStatus{}, nil
	case MsgTypeUpdateAvailable:
		return &UpdateAvailable{}, nil
	case MsgTypeWorkflowRun:
		return &WorkflowRun{}, nil
	case MsgTypeWorkflowStop:
		return &WorkflowStop{}, nil
	case MsgTypeWorkflowDesignChat:
		return &WorkflowDesignChat{}, nil
	case MsgTypeWorkflowStarted:
		return &WorkflowStarted{}, nil
	case MsgTypeWorkflowNodeStarted:
		return &WorkflowNodeStarted{}, nil
	case MsgTypeWorkflowNodeStream:
		return &WorkflowNodeStream{}, nil
	case MsgTypeWorkflowNodeComplete:
		return &WorkflowNodeComplete{}, nil
	case MsgTypeWorkflowNodeError:
		return &WorkflowNodeError{}, nil
	case MsgTypeWorkflowComplete:
		return &WorkflowComplete{}, nil
	case MsgTypeWorkflowError:
		return &WorkflowError{}, nil
	case MsgTypeWorkflowDesignChatStream:
		return &WorkflowDesignChatStream{}, nil
	case MsgTypeWorkflowDesignChatComplete:
		return &WorkflowDesignChatComplete{}, nil
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
