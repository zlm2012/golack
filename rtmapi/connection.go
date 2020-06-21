package rtmapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/oklahomer/golack/event"
	"github.com/tidwall/gjson"
	"io"
)

type UnexpectedMessageTypeError struct {
	MessageType int
	Payload     []byte
}

func (e *UnexpectedMessageTypeError) Error() string {
	return fmt.Sprintf("unexpected message type, %d, is given: %s", e.MessageType, e.Payload)
}

type DecodedPayload interface{}

type PayloadReceiver interface {
	Receive() (DecodedPayload, error)
}

type PayloadSender interface {
	Send(*OutgoingMessage) error
	Ping() error
}

type Connection interface {
	PayloadReceiver
	PayloadSender
	io.Closer
}

// Connect connects to Slack WebSocket server.
func Connect(_ context.Context, url string) (Connection, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
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
	messageType, payload, err := wrapper.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	// Only TextMessage is supported by RTM API.
	if messageType != websocket.TextMessage {
		return nil, &UnexpectedMessageTypeError{MessageType: messageType, Payload: payload}
	}

	decoded, err := decodePayload(payload)
	return decoded, err
}

func (wrapper *connWrapper) Send(message *OutgoingMessage) error {
	// ID must be unique per connection.
	// Manage this value at here.
	message.ID = wrapper.outgoingEventID.Next()
	return wrapper.conn.WriteJSON(message)
}

func (wrapper *connWrapper) Ping() error {
	ping := NewPing(wrapper.outgoingEventID)
	return wrapper.conn.WriteJSON(ping)
}

func (wrapper *connWrapper) Close() error {
	return wrapper.conn.Close()
}

func decodePayload(input json.RawMessage) (DecodedPayload, error) {
	// Sometimes an empty payload comes in.
	// Return a designated error and let caller decide how to handle.
	input = bytes.TrimSpace(input)
	if len(input) == 0 {
		return nil, event.ErrEmptyPayload
	}

	// Map the payload to predefined event
	parsed := gjson.ParseBytes(input)
	e, err := event.Map(parsed)
	if err == nil {
		return e, nil
	}

	// Error may be returned when a WebSocket protocol-specific payload is given because such payload is not listed as "event" on https://api.slack.com/events.

	// Type is not defined for "event," but can be for WebSocket protocol
	if parsed.Get("reply_to").Exists() {
		// https://api.slack.com/rtm#ping_and_pong
		payloadType := parsed.Get("type")
		if payloadType.Exists() && payloadType.String() == "pong" {
			return decodePong(input)
		}

		// https://api.slack.com/rtm#handling_responses
		payloadOK := parsed.Get("ok")
		if payloadOK.Exists() {
			if payloadOK.Bool() {
				return decodeOKResponse(input)
			}
			return decodeNGResponse(input)
		}
	}

	return nil, event.NewMalformedPayloadError(fmt.Sprintf("given json object has unknown structure. can not handle: %s.", input))
}

func decodePong(input json.RawMessage) (DecodedPayload, error) {
	mapping := &Pong{}
	err := json.Unmarshal(input, mapping)
	if err != nil {
		return nil, event.NewMalformedPayloadError(fmt.Sprintf("malformed pong payload is given: %s", input))
	}
	return mapping, nil
}

func decodeOKResponse(input json.RawMessage) (DecodedPayload, error) {
	mapping := &OKReply{}
	err := json.Unmarshal(input, mapping)
	if err != nil {
		return nil, event.NewMalformedPayloadError(fmt.Sprintf("malformed response is given: %s", input))
	}
	return mapping, nil
}

func decodeNGResponse(input json.RawMessage) (DecodedPayload, error) {
	mapping := &NGReply{}
	err := json.Unmarshal(input, mapping)
	if err != nil {
		return nil, event.NewMalformedPayloadError(fmt.Sprintf("malformed response is given: %s", input))
	}
	return mapping, nil
}
