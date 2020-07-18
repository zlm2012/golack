package rtmapi

import (
	"github.com/oklahomer/golack/v2/event"
)

// Pong is given when client send Ping.
// https://api.slack.com/rtm#ping_and_pong
type Pong struct {
	event.TypedEvent
	ReplyTo uint `json:"reply_to"`
}
