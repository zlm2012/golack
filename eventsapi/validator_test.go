package eventsapi

import (
	"testing"
	"time"
)

func TestSignatureValidator_Validate(t *testing.T) {
	request := &SlackRequest{
		Signature: "v0=fd5b65c1b8655daf297b59df9156cc113b3d1f705dff63356888ccca73ef91fb",
		TimeStamp: time.Unix(123456789, 0),
		Payload:   []byte("command=/weather&text=94070"),
	}

	validator := &SignatureValidator{Secret: "Shhh"}
	valid := validator.Validate(request)

	if !valid {
		t.Fatal("Validator unexpectedly failed")
	}
}
