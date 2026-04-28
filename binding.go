package protocol

import (
	"fmt"
	"reflect"
)

// Binding contains the protocol identifiers used to bind requests, sessions,
// channels, and transport streams.
type Binding struct {
	RequestID  string
	SessionID  string
	ChannelID  string
	MachineID  string
	PreviewID  string
	StreamID   string
	TerminalID string
	TaskID     string
}

// ExtractBinding returns common correlation identifiers from a protocol
// message payload or Envelope.
func ExtractBinding(message any) Binding {
	if env, ok := message.(*Envelope); ok && env != nil {
		message = env.Payload
	}

	value := reflect.ValueOf(message)
	if !value.IsValid() {
		return Binding{}
	}
	for value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return Binding{}
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return Binding{}
	}

	return Binding{
		RequestID:  stringField(value, "RequestID"),
		SessionID:  stringField(value, "SessionID"),
		ChannelID:  stringField(value, "ChannelID"),
		MachineID:  stringField(value, "MachineID"),
		PreviewID:  stringField(value, "PreviewID"),
		StreamID:   stringField(value, "StreamID"),
		TerminalID: stringField(value, "TerminalID"),
		TaskID:     stringField(value, "TaskID"),
	}
}

// ValidateRequestBinding verifies that an actual message is bound to the same
// request ID as the expected request message or binding.
func ValidateRequestBinding(expected any, actual any) error {
	want := bindingFromAny(expected)
	got := bindingFromAny(actual)
	if want.RequestID == "" {
		return fmt.Errorf("binding: expected requestId is empty")
	}
	if got.RequestID == "" {
		return fmt.Errorf("binding: actual requestId is empty")
	}
	if got.RequestID != want.RequestID {
		return fmt.Errorf("binding: requestId mismatch: got %q, want %q", got.RequestID, want.RequestID)
	}
	return nil
}

// ValidateSessionBinding verifies that an actual message is bound to the same
// session and channel identifiers as the expected message or binding.
func ValidateSessionBinding(expected any, actual any) error {
	want := bindingFromAny(expected)
	got := bindingFromAny(actual)
	if want.SessionID == "" {
		return fmt.Errorf("binding: expected sessionId is empty")
	}
	if got.SessionID == "" {
		return fmt.Errorf("binding: actual sessionId is empty")
	}
	if got.SessionID != want.SessionID {
		return fmt.Errorf("binding: sessionId mismatch: got %q, want %q", got.SessionID, want.SessionID)
	}

	if want.ChannelID != "" {
		if got.ChannelID == "" {
			return fmt.Errorf("binding: actual channelId is empty")
		}
		if got.ChannelID != want.ChannelID {
			return fmt.Errorf("binding: channelId mismatch: got %q, want %q", got.ChannelID, want.ChannelID)
		}
	}
	return nil
}

func bindingFromAny(value any) Binding {
	if binding, ok := value.(Binding); ok {
		return binding
	}
	if binding, ok := value.(*Binding); ok && binding != nil {
		return *binding
	}
	return ExtractBinding(value)
}

func stringField(value reflect.Value, name string) string {
	field := value.FieldByName(name)
	if !field.IsValid() || field.Kind() != reflect.String {
		return ""
	}
	return field.String()
}
