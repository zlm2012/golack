package rtmapi

import (
	"encoding/json"
	"github.com/oklahomer/golack/v2/event"
	"strings"
	"testing"
)

func TestMarshalPingEvent(t *testing.T) {
	ping := &Ping{
		OutgoingEvent: OutgoingEvent{
			ID:         1,
			TypedEvent: event.TypedEvent{Type: "ping"},
		},
	}
	val, err := json.Marshal(ping)
	if err != nil {
		t.Fatalf("error occured while encoding. %s.", err.Error())
	}

	if strings.Contains(string(val), "ping") != true {
		t.Fatalf(`returned string doesn't contain "ping". %s.`, string(val))
	}
}

func TestUnmarshalPingEvent(t *testing.T) {
	str := `{"type": "ping", "id": 123}`
	ping := &Ping{}
	if err := json.Unmarshal([]byte(str), ping); err != nil {
		t.Errorf("error on Unmarshal. %s.", err.Error())
		return
	}

	if ping.Type != "ping" {
		t.Errorf("something is wrong with unmarshaled result. %#v.", ping)
	}

	if ping.ID != 123 {
		t.Errorf("unmarshaled id is wrong %d. expecting %d.", ping.ID, 123)
	}
}

func TestNewOutgoingMessage(t *testing.T) {
	channelID := event.ChannelID("channel")
	message := "dummy message"

	outgoingMessage := NewOutgoingMessage(channelID, message)

	if outgoingMessage.ChannelID != channelID {
		t.Errorf("Passed channelID is not set: %s.", outgoingMessage.ChannelID)
	}

	if outgoingMessage.Text != message {
		t.Errorf("Passed message is not set: %s.", outgoingMessage.Text)
	}
}

func TestOutgoingMessage_WithThreadTimeStamp(t *testing.T) {
	timeStamp := &event.TimeStamp{}
	message := &OutgoingMessage{}

	message.WithThreadTimeStamp(timeStamp)

	if message.ThreadTimeStamp != timeStamp {
		t.Errorf("Passed timestamp is not set: %s.", message.ThreadTimeStamp)
	}
}
