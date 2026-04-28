package protocol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

const (
	DefaultMaxFrameBytes   = 1 << 20
	DefaultMaxDepth        = 32
	DefaultMaxObjectFields = 256
	DefaultMaxArrayItems   = 4096
)

// EnvelopeLimits bounds protocol frame parsing before payload unmarshalling.
// Zero values use the default limit for that field.
type EnvelopeLimits struct {
	MaxFrameBytes   int
	MaxDepth        int
	MaxObjectFields int
	MaxArrayItems   int
}

// DefaultEnvelopeLimits returns the default frame parsing limits.
func DefaultEnvelopeLimits() EnvelopeLimits {
	return EnvelopeLimits{
		MaxFrameBytes:   DefaultMaxFrameBytes,
		MaxDepth:        DefaultMaxDepth,
		MaxObjectFields: DefaultMaxObjectFields,
		MaxArrayItems:   DefaultMaxArrayItems,
	}
}

// ValidateEnvelopeFrame checks transport-level JSON bounds before a frame is
// unmarshaled into a concrete protocol payload.
func ValidateEnvelopeFrame(data []byte, limits EnvelopeLimits) error {
	limits = normalizeEnvelopeLimits(limits)
	if limits.MaxFrameBytes > 0 && len(data) > limits.MaxFrameBytes {
		return fmt.Errorf("frame exceeds max size: %d > %d bytes", len(data), limits.MaxFrameBytes)
	}

	validator := frameValidator{
		limits: limits,
		dec:    json.NewDecoder(bytes.NewReader(data)),
	}
	validator.dec.UseNumber()
	return validator.validate()
}

type frameValidator struct {
	limits   EnvelopeLimits
	dec      *json.Decoder
	stack    []jsonContainer
	rootSeen bool
	rootKind byte
}

type jsonContainer struct {
	kind         byte
	count        int
	expectingKey bool
}

func normalizeEnvelopeLimits(limits EnvelopeLimits) EnvelopeLimits {
	defaults := DefaultEnvelopeLimits()
	if limits.MaxFrameBytes == 0 {
		limits.MaxFrameBytes = defaults.MaxFrameBytes
	}
	if limits.MaxDepth == 0 {
		limits.MaxDepth = defaults.MaxDepth
	}
	if limits.MaxObjectFields == 0 {
		limits.MaxObjectFields = defaults.MaxObjectFields
	}
	if limits.MaxArrayItems == 0 {
		limits.MaxArrayItems = defaults.MaxArrayItems
	}
	return limits
}

func (v *frameValidator) validate() error {
	for {
		tok, err := v.dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("frame json: %w", err)
		}
		if err := v.consume(tok); err != nil {
			return err
		}
	}
	if !v.rootSeen {
		return fmt.Errorf("frame json: empty")
	}
	if v.rootKind != '{' {
		return fmt.Errorf("frame root must be object")
	}
	if len(v.stack) != 0 {
		return fmt.Errorf("frame json: incomplete document")
	}
	return nil
}

func (v *frameValidator) consume(tok any) error {
	switch value := tok.(type) {
	case json.Delim:
		return v.consumeDelim(value)
	case string:
		return v.consumeString()
	default:
		return v.markValue(0)
	}
}

func (v *frameValidator) consumeDelim(delim json.Delim) error {
	switch delim {
	case '{':
		if err := v.markValue('{'); err != nil {
			return err
		}
		if len(v.stack)+1 > v.limits.MaxDepth {
			return fmt.Errorf("frame exceeds max depth: %d > %d", len(v.stack)+1, v.limits.MaxDepth)
		}
		v.stack = append(v.stack, jsonContainer{kind: '{', expectingKey: true})
		return nil
	case '[':
		if err := v.markValue('['); err != nil {
			return err
		}
		if len(v.stack)+1 > v.limits.MaxDepth {
			return fmt.Errorf("frame exceeds max depth: %d > %d", len(v.stack)+1, v.limits.MaxDepth)
		}
		v.stack = append(v.stack, jsonContainer{kind: '[', expectingKey: false})
		return nil
	case '}':
		return v.closeContainer('{')
	case ']':
		return v.closeContainer('[')
	default:
		return fmt.Errorf("frame json: invalid delimiter %q", delim)
	}
}

func (v *frameValidator) consumeString() error {
	if len(v.stack) > 0 {
		top := &v.stack[len(v.stack)-1]
		if top.kind == '{' && top.expectingKey {
			top.count++
			if v.limits.MaxObjectFields > 0 && top.count > v.limits.MaxObjectFields {
				return fmt.Errorf("frame object exceeds max fields: %d > %d", top.count, v.limits.MaxObjectFields)
			}
			top.expectingKey = false
			return nil
		}
	}
	return v.markValue(0)
}

func (v *frameValidator) markValue(kind byte) error {
	if len(v.stack) == 0 {
		if v.rootSeen {
			return fmt.Errorf("frame json: multiple top-level values")
		}
		v.rootSeen = true
		v.rootKind = kind
		return nil
	}

	top := &v.stack[len(v.stack)-1]
	switch top.kind {
	case '{':
		if top.expectingKey {
			return fmt.Errorf("frame json: expected object key")
		}
		top.expectingKey = true
	case '[':
		top.count++
		if v.limits.MaxArrayItems > 0 && top.count > v.limits.MaxArrayItems {
			return fmt.Errorf("frame array exceeds max items: %d > %d", top.count, v.limits.MaxArrayItems)
		}
	}
	return nil
}

func (v *frameValidator) closeContainer(kind byte) error {
	if len(v.stack) == 0 {
		return fmt.Errorf("frame json: unexpected close delimiter")
	}
	top := v.stack[len(v.stack)-1]
	if top.kind != kind {
		return fmt.Errorf("frame json: mismatched close delimiter")
	}
	if top.kind == '{' && !top.expectingKey {
		return fmt.Errorf("frame json: missing object value")
	}
	v.stack = v.stack[:len(v.stack)-1]
	return nil
}
