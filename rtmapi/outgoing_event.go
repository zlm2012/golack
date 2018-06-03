package rtmapi

import "github.com/oklahomer/golack/slackobject"

// OutgoingEvent takes care of some common fields all outgoing event MUST have.
// https://api.slack.com/rtm#events
type OutgoingEvent struct {
	TypedEvent

	// ID is an unique identifier that client declares.
	// https://api.slack.com/rtm#sending_messages
	// Every event should have a unique (for that connection) positive integer ID.
	// All replies to that message will include this ID allowing the client to correlate responses with the messages sent;
	// replies may be "out of order" due to the asynchronous nature of the message servers.
	ID uint `json:"id"`
}

// OutgoingMessage represents a simple message sent from client to Slack server via WebSocket connection.
// This is the only format RTM API supports. To send more richly formatted message, use Web API.
// https://api.slack.com/rtm#sending_messages
type OutgoingMessage struct {
	OutgoingEvent
	Channel slackobject.ChannelID `json:"channel"`
	Text    string                `json:"text"`
}

// NewOutgoingMessage is a constructor to create new OutgoingMessage instance with given arguments.
func NewOutgoingMessage(eventID *OutgoingEventID, channel slackobject.ChannelID, text string) *OutgoingMessage {
	return &OutgoingMessage{
		Channel: channel,
		Text:    text,
		OutgoingEvent: OutgoingEvent{
			ID: eventID.Next(),
			TypedEvent: TypedEvent{
				Type: MessageEvent,
			},
		},
	}
}

// Ping is an event that can be sent to slack endpoint via WebSocket to see if the connection is alive.
// Slack sends back Pong event if connection is O.K.
type Ping struct {
	OutgoingEvent
}

// NewPing creates new Ping instance with given arguments.
func NewPing(eventID *OutgoingEventID) *Ping {
	return &Ping{
		OutgoingEvent: OutgoingEvent{
			ID:         eventID.Next(),
			TypedEvent: TypedEvent{Type: PingEvent},
		},
	}
}
