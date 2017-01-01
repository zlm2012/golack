package rtmapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
	"io"
	"strings"
)

var (
	ErrEmptyPayload = errors.New("empty payload was given")
)

type DecodedPayload interface{}

type PayloadReceiver interface {
	Receive() (DecodedPayload, error)
}

type PayloadSender interface {
	Send(*Channel, string) error
	Ping() error
}

type Connection interface {
	PayloadReceiver
	PayloadSender
	io.Closer
}

// Connect connects to Slack WebSocket server.
func Connect(_ context.Context, url string) (Connection, error) {
	conn, err := websocket.Dial(url, "", "http://localhost")
	if err != nil {
		return nil, err
	}

	return newConnectionWrapper(conn), nil
}

// connWrapper is a thin wrapper that wraps WebSocket connection and its methods.
// This instance is created per-connection.
type connWrapper struct {
	conn *websocket.Conn

	// https://api.slack.com/rtm#sending_messages
	// Every event should have a unique (for that connection) positive integer ID.
	outgoingEventID *OutgoingEventID
}

func newConnectionWrapper(conn *websocket.Conn) Connection {
	return &connWrapper{
		conn:            conn,
		outgoingEventID: NewOutgoingEventID(),
	}

}

// Receive is a blocking method to receive payload from WebSocket connection.
// When connection is closed in the middle of this method call, this immediately returns error.
func (wrapper *connWrapper) Receive() (DecodedPayload, error) {
	// Slack's RTM events and reply all have different form
	// so websocket.JSON.Receive can only work with json.RawMessage,
	// which ends up with multiple json.Unmarshal calls for proper mapping.
	payload := []byte{}
	err := websocket.Message.Receive(wrapper.conn, &payload)
	if err != nil {
		return nil, err
	}

	return decodePayload(payload)
}

func (wrapper *connWrapper) Send(channel *Channel, content string) error {
	event := NewOutgoingMessage(wrapper.outgoingEventID, channel, content)
	return websocket.JSON.Send(wrapper.conn, event)
}

func (wrapper *connWrapper) Ping() error {
	ping := NewPing(wrapper.outgoingEventID)
	return websocket.JSON.Send(wrapper.conn, ping)
}

func (wrapper *connWrapper) Close() error {
	return wrapper.conn.Close()
}

func decodePayload(input json.RawMessage) (DecodedPayload, error) {
	inputStr := strings.TrimSpace(string(input))
	if len(inputStr) == 0 {
		return nil, ErrEmptyPayload
	}

	res := gjson.Parse(inputStr)
	eventType := res.Get("type")
	if eventType.Exists() {
		typeVal := eventType.String()
		switch typeVal {
		case HelloEvent:
			return unmarshal(input, &Hello{})
		case MessageEvent:
			subType := res.Get("subtype")
			if subType.Exists() {
				switch subType.String() {
				// TODO handle each subtypes
				case BotMessage:
					return unmarshal(input, &MiscMessage{})
				case ChannelArchive:
					return unmarshal(input, &MiscMessage{})
				case ChannelJoin:
					return unmarshal(input, &MiscMessage{})
				case ChannelLeave:
					return unmarshal(input, &MiscMessage{})
				case ChannelName:
					return unmarshal(input, &MiscMessage{})
				case ChannelPurpose:
					return unmarshal(input, &MiscMessage{})
				case ChannelTopic:
					return unmarshal(input, &MiscMessage{})
				case ChannelUnarchive:
					return unmarshal(input, &MiscMessage{})
				case FileComment:
					return unmarshal(input, &MiscMessage{})
				case FileMention:
					return unmarshal(input, &MiscMessage{})
				case FileShare:
					return unmarshal(input, &MiscMessage{})
				case GroupArchive:
					return unmarshal(input, &MiscMessage{})
				case GroupJoin:
					return unmarshal(input, &MiscMessage{})
				case GroupLeave:
					return unmarshal(input, &MiscMessage{})
				case GroupName:
					return unmarshal(input, &MiscMessage{})
				case GroupPurpose:
					return unmarshal(input, &MiscMessage{})
				case GroupTopic:
					return unmarshal(input, &MiscMessage{})
				case GroupUnarchive:
					return unmarshal(input, &MiscMessage{})
				case MeMessage:
					return unmarshal(input, &MiscMessage{})
				case MessageChanged:
					return unmarshal(input, &MiscMessage{})
				case MessageDeleted:
					return unmarshal(input, &MiscMessage{})
				case PinnedItem:
					return unmarshal(input, &MiscMessage{})
				case UnpinnedItem:
					return unmarshal(input, &MiscMessage{})
				default:
					return unmarshal(input, &MiscMessage{})
				}
			} else {
				// plain message
				return unmarshal(input, &Message{})
			}
		case TeamMigrationStartedEvent:
			return unmarshal(input, &TeamMigrationStarted{})
		case PongEvent:
			return unmarshal(input, &Pong{})
		default:
			return nil, NewMalformedPayloadError(fmt.Sprintf(`unsupported event type "%s" is given: %s`, typeVal, input))
		}
	}

	if res.Get("reply_to").Exists() && res.Get("ok").Exists() {
		// https://api.slack.com/rtm#handling_responses
		// When incoming object can't be treated as Event, try treat this as WebSocketReply.
		return unmarshal(input, &WebSocketReply{})
	}

	return nil, NewMalformedPayloadError(fmt.Sprintf("given json object has unknown structure. can not handle: %s.", input))
}

func unmarshal(input json.RawMessage, mapping interface{}) (DecodedPayload, error) {
	if err := json.Unmarshal(input, mapping); err != nil {
		return nil, NewMalformedPayloadError(err.Error())
	}

	return mapping, nil
}
