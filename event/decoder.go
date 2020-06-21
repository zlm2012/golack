package event

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"reflect"
)

var (
	ErrEmptyPayload = errors.New("empty payload was given")
)

// MalformedPayloadError represents an error that given JSON payload is not properly formatted.
// e.g. required fields are not given in an expected format, or payload is not a valid JSON string.
type MalformedPayloadError struct {
	Err string
}

// Error returns its error string.
func (e *MalformedPayloadError) Error() string {
	return e.Err
}

// NewMalformedPayloadError creates new MalformedPayloadError instance with given arguments.
func NewMalformedPayloadError(str string) *MalformedPayloadError {
	return &MalformedPayloadError{Err: str}
}

type UnknownPayloadTypeError struct {
	Err string
}

func (e *UnknownPayloadTypeError) Error() string {
	return e.Err
}

func NewUnknownPayloadTypeError(str string) *UnknownPayloadTypeError {
	return &UnknownPayloadTypeError{Err: str}
}

func Decode(payload []byte) (interface{}, error) {
	payload = bytes.TrimSpace(payload)
	if len(payload) == 0 {
		return nil, ErrEmptyPayload
	}

	// See if required "type" field exists
	parsed := gjson.ParseBytes(payload)
	return Map(parsed)
}

func Map(parsed gjson.Result) (interface{}, error) {
	typeObject := parsed.Get("type")
	if !typeObject.Exists() || typeObject.Type != gjson.String {
		return nil, NewMalformedPayloadError(fmt.Sprintf("given payload has unknown structure. can not handle: %s", parsed))
	}

	// Handle those events that requires extra care
	typeValue := typeObject.String()
	if typeValue == "message" {
		// A message event may have subtype field
		subtypeValue := parsed.Get("subtype")
		if subtypeValue.Exists() {
			mapping, ok := subTypeMap[subtypeValue.String()]
			if !ok {
				return nil, NewUnknownPayloadTypeError(fmt.Sprintf("unknown subtype of %s is given: %s", subtypeValue, parsed))
			}
			return unmarshal([]byte(parsed.String()), mapping)
		}

		channelType := parsed.Get("channel_type")
		if channelType.Exists() {
			mapping, ok := messageChannelTypeMap[channelType.String()]
			if !ok {
				return nil, NewUnknownPayloadTypeError(fmt.Sprintf("unknown channel_type of %s is given: %s", channelType, parsed))
			}
			return unmarshal([]byte(parsed.String()), mapping)
		}
	}

	// Map to the corresponding struct
	mapping, ok := eventTypeMap[typeValue]
	if !ok {
		return nil, NewUnknownPayloadTypeError(fmt.Sprintf("unknown type of %s is given: %s", typeValue, parsed))
	}

	return unmarshal([]byte(parsed.String()), mapping)
}

func unmarshal(input []byte, mapping reflect.Type) (interface{}, error) {
	payload := reflect.New(mapping).Interface()
	err := json.Unmarshal(input, &payload)
	if err != nil {
		return nil, NewMalformedPayloadError(err.Error())
	}

	return payload, nil
}
