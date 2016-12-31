package rtmapi

import (
	"bytes"
	"testing"
)

func TestChannel_UnmarshalText(t *testing.T) {
	name := "myChannel"
	channel := &Channel{}
	if err := channel.UnmarshalText([]byte(name)); err != nil {
		t.Fatalf("error on channel unmarshal: %s.", err.Error())
	}

	if channel.Name != name {
		t.Errorf("expected value is not returned: %s.", channel.Name)
	}
}

func TestChannel_MarshalText(t *testing.T) {
	name := "myChannel"
	channel := &Channel{Name: name}

	b, err := channel.MarshalText()

	if err != nil {
		t.Fatalf("unexpected error on channel marshal: %s.", err.Error())
	}

	if !bytes.Equal(b, []byte(name)) {
		t.Errorf("marshaled value is wrong: %s. expected: %s.", b, name)
	}
}
