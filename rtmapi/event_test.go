package rtmapi

import (
	"bytes"
	"strings"
	"testing"
)

func TestEventType_UnmarshalText(t *testing.T) {
	for _, given := range possibleEvents {
		var eventType EventType
		if err := eventType.UnmarshalText([]byte(given)); err != nil {
			t.Errorf("Unexpected error on valid EventType unmarshal: %s.", err.Error())
			continue
		}

		if strings.Compare(string(eventType), string(given)) != 0 {
			t.Errorf("Unexpected EventType was returned: %s. Expected: %s.", eventType, given)
		}
	}
}

func TestEventType_UnmarshalText_UnsupportedEvent(t *testing.T) {
	var eventType EventType
	if err := eventType.UnmarshalText([]byte("UNKNOWN")); err != nil {
		t.Errorf("Unexpected error on EventTyep unmarshal: %s.", err.Error())
		return
	}

	if strings.Compare(string(eventType), string(UnsupportedEvent)) != 0 {
		t.Errorf("Unexpected EventType was returned: %s. Expected: %s.", eventType, UnsupportedEvent)
	}
}

func TestEventType_MarshalText(t *testing.T) {
	for _, given := range possibleEvents {
		text, err := given.MarshalText()
		if err != nil {
			t.Errorf("Unexpected error on EventType marshal: %s.", err.Error())
			continue
		}

		if !bytes.Equal(text, []byte(given)) {
			t.Errorf("Unexpected value is returned: %s. Expected: %s.", string(text), string(given))
		}
	}
}

func TestEventType_MarshalText_ZeroValue(t *testing.T) {
	var eventType EventType
	text, err := eventType.MarshalText()

	if err != nil {
		t.Fatalf("Unexpected error on EventType marshal: %s.", err.Error())
	}

	if !bytes.Equal(text, []byte(UnsupportedEvent)) {
		t.Errorf("Unexpected value is returned: %s. Expected: %s.", string(text), string(UnsupportedEvent))
	}
}

func TestSubType_UnmarshalText(t *testing.T) {
	for _, given := range possibleSubTypes {
		var subType SubType
		if err := subType.UnmarshalText([]byte(given)); err != nil {
			t.Errorf("Unexpected error on SubType unmarshal: %s.", err.Error())
			continue
		}

		if strings.Compare(string(subType), string(given)) != 0 {
			t.Errorf("Unexpected value is returned: %s. Expected: %s.", subType, given)
		}
	}
}

func TestSubType_UnmarshalText_Empty(t *testing.T) {
	var subType SubType
	if err := subType.UnmarshalText([]byte("")); err != nil {
		t.Fatalf("Unexpected error on SubType unmarshal: %s.", err.Error())
	}

	if strings.Compare(string(subType), string(Empty)) != 0 {
		t.Errorf("Unexpected value is returned: %s. Expected: %s.", subType, Empty)
	}
}

func TestSubType_MarshalText(t *testing.T) {
	for _, given := range possibleSubTypes {
		text, err := given.MarshalText()
		if err != nil {
			t.Errorf("Unexpected error on SubType marshal: %s.", err.Error())
			continue
		}

		if !bytes.Equal(text, []byte(given)) {
			t.Errorf("Unexpected value is returned: %s. Expected: %s.", string(text), string(given))
		}
	}
}

func TestSubType_MarshalText_ZeroValue(t *testing.T) {
	var subType SubType
	text, err := subType.MarshalText()
	if err != nil {
		t.Fatalf("Unexpected error on SubType marshal: %s.", err.Error())
	}

	if !bytes.Equal(text, []byte("")) {
		t.Errorf("Unexpected value is returned: %s. Expected: %s.", string(text), string(""))
	}
}
