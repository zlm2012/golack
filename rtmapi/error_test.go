package rtmapi

import (
	"github.com/oklahomer/golack/event"
	"testing"
)

func TestNewMalformedPayloadError(t *testing.T) {
	str := "myError"
	err := event.NewMalformedPayloadError(str)

	if err == nil {
		t.Fatal("error instance is not returned.")
	}

	if err.Err != str {
		t.Errorf("expected error string is not set: %s.", err.Error())
	}
}

func TestMalformedPayloadError_Error(t *testing.T) {
	str := "myErr"
	err := event.MalformedPayloadError{Err: str}

	if err.Error() != str {
		t.Errorf("expected error string is not returned: %s.", err.Error())
	}
}
