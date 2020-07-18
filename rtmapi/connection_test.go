package rtmapi

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/oklahomer/golack/v2/event"
	"github.com/oklahomer/golack/v2/testutil"
	"net"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestConnect(t *testing.T) {
	testutil.RunWithWebSocket(func(addr net.Addr) {
		url := fmt.Sprintf("ws://%s%s", addr, "/echo")
		connection, err := Connect(context.TODO(), url)
		if err != nil {
			t.Fatalf("webSocket connection error: %s.", err.Error())
		}

		if connection == nil {
			t.Fatalf("Connection is not reurned.")
		}
		connection.Close()
	})
}

func TestConnect_Fail(t *testing.T) {
	testutil.RunWithWebSocket(func(addr net.Addr) {
		// Establish connection
		url := fmt.Sprintf("ws://%s%s", addr, "/undefined_path")
		_, err := Connect(context.TODO(), url)

		if err == nil {
			t.Fatal("expected error is not returned.")
		}

		if err != websocket.ErrBadHandshake {
			t.Fatalf("Unexpected error is returned: %s.", err.Error())
		}
	})
}

func Test_newConnectionWrapper(t *testing.T) {
	conn := newConnectionWrapper(&websocket.Conn{})

	if conn == nil {
		t.Fatal("connection is not returned.")
	}

	wrapper, ok := conn.(*connWrapper)
	if !ok {
		t.Errorf("unexpected type is represented: %#v.", conn)
	}

	if wrapper.conn == nil {
		t.Error("websocket.Conn is not wrapped.")
	}

	if wrapper.outgoingEventID == nil {
		t.Error("id dispenser is not initialized and set.")
	}
}

func TestConnWrapper_Receive(t *testing.T) {
	testutil.RunWithWebSocket(func(addr net.Addr) {
		url := fmt.Sprintf("ws://%s%s", addr, "/echo")
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			t.Fatal("can't establish connection with test server")
		}
		defer conn.Close()

		var channelID event.ChannelID = "C12345"
		var userID event.UserID = "U6789"
		text := "Hello world!"
		timestamp := 1355517523
		slackTimestamp := fmt.Sprintf("%d.000005", timestamp)
		input := fmt.Sprintf(`{"type": "message", "channel": "%s", "user": "%s", "text": "%s", "ts": "%s"}`, channelID.String(), userID.String(), text, slackTimestamp)

		connWrapper := newConnectionWrapper(conn)
		conn.WriteMessage(websocket.TextMessage, []byte(input))
		decodedPayload, err := connWrapper.Receive()
		if err != nil {
			t.Fatalf("error on payload reception: %s.", err.Error())
		}

		message, ok := decodedPayload.(*event.Message)
		if !ok {
			t.Fatalf("received payload is not Message: %#v.", message)
		}

		if message.ChannelID != channelID {
			t.Errorf("expected channel name is not given: %s.", message.ChannelID.String())
		}

		if message.SenderID != userID {
			t.Errorf("expected user is not given: %s.", message.SenderID)
		}

		if message.Text != text {
			t.Errorf("expected text is not given: %s.", message.Text)
		}

		if !message.TimeStamp.Time.Equal(time.Unix(1355517523, 0)) {
			t.Errorf("expected time is not given: %d.", message.TimeStamp.Time.Unix())
		}

		if message.TimeStamp.OriginalValue != slackTimestamp {
			t.Errorf("expected timestamp original value is not given: %s.", message.TimeStamp.OriginalValue)
		}

	})
}

func TestConnWrapper_Send(t *testing.T) {
	testutil.RunWithWebSocket(func(addr net.Addr) {
		url := fmt.Sprintf("ws://%s%s", addr, "/echo")
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			t.Fatal("can't establish connection with test server")
		}
		defer conn.Close()

		connWrapper := newConnectionWrapper(conn)
		message := &OutgoingMessage{}
		if err := connWrapper.Send(message); err != nil {
			t.Errorf("error on sending message over WebSocket connection. %#v.", err)
		}

		if message.ID == 0 {
			t.Errorf("Send() method must append message id.")
		}
	})
}

func TestConnWrapper_Ping(t *testing.T) {
	testutil.RunWithWebSocket(func(addr net.Addr) {
		url := fmt.Sprintf("ws://%s%s", addr, "/ping")
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			t.Fatal("can't establish connection with test server")
		}
		defer conn.Close()

		connWrapper := newConnectionWrapper(conn)
		if err := connWrapper.Ping(); err != nil {
			t.Errorf("error on sending message over WebSocket connection. %#v.", err)
		}
	})
}

func TestConnWrapper_Close(t *testing.T) {
	testutil.RunWithWebSocket(func(addr net.Addr) {
		url := fmt.Sprintf("ws://%s%s", addr, "/echo")
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			t.Fatal("can't establish connection with test server")
		}

		connWrapper := newConnectionWrapper(conn)

		if err := connWrapper.Close(); err != nil {
			t.Fatal("error on connection close")
		}

		if err := conn.Close(); err == nil {
			t.Fatal("net.OpError should be returned when WebSocket.Conn.Close is called multiple times.")
		}
	})
}

func Test_decodePayload(t *testing.T) {
	type expected struct {
		value interface{}
		err   interface{}
	}
	tests := []struct {
		input    string
		expected *expected
	}{
		{
			input: `{"reply_to": 1234, "type": "pong", "time": 1403299273342}`,
			expected: &expected{
				value: &Pong{
					TypedEvent: event.TypedEvent{
						Type: "pong",
					},
					ReplyTo: 1234,
				},
			},
		},
		{
			input: `{"ok": true, "reply_to": 1, "ts": "1355517523.000005", "text": "Hello world"}`,
			expected: &expected{
				value: &OKReply{
					Reply: Reply{
						OK:      true,
						ReplyTo: 1,
					},
					TimeStamp: &event.TimeStamp{
						Time:          time.Unix(1355517523, 0),
						OriginalValue: "1355517523.000005",
					},
					Text: "Hello world",
				},
			},
		},
		{
			input: `{"ok": false, "reply_to": 1, "error": {"code": 2, "msg": "message text is missing"}}`,
			expected: &expected{
				value: &NGReply{
					Reply: Reply{
						OK:      false,
						ReplyTo: 1,
					},
					Error: ReplyErrorReason{
						Code:    2,
						Message: "message text is missing",
					},
				},
			},
		},
		{
			// invalid type
			input: `{"type": "unsupportedEventType"}`,
			expected: &expected{
				err: reflect.TypeOf(&event.MalformedPayloadError{}),
			},
		},
		{
			// required fields are not given
			input: `{"what": true}`,
			expected: &expected{
				value: nil,
				err:   reflect.TypeOf(&event.MalformedPayloadError{}),
			},
		},
		{
			// not valid json structure
			input: `malformedJson`,
			expected: &expected{
				err: reflect.TypeOf(&event.MalformedPayloadError{}),
			},
		},
		{
			input: "　",
			expected: &expected{
				err: event.ErrEmptyPayload,
			},
		},
		{
			input: "\r",
			expected: &expected{
				err: event.ErrEmptyPayload,
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			payload, err := decodePayload([]byte(tt.input))

			if tt.expected.value != nil {
				if err != nil {
					t.Fatalf("Unexpected error is returned: %#v", err)
				}

				if !reflect.DeepEqual(payload, tt.expected.value) {
					t.Fatalf("Expected %#v but was %#v", tt.expected.value, payload)
				}

				return
			}

			expectedErr := tt.expected.err
			// See if pre-declared error instance is returned
			if reflect.TypeOf(expectedErr) == reflect.TypeOf(errors.New("DUMMY")) {
				// pre-declared error instance is returned.
				if expectedErr != err {
					t.Fatalf("Expected error is not returned: %#v", err)
				}
				return
			}

			// Error type comparison is required
			if expectedErr != reflect.TypeOf(err) {
				t.Fatalf("Unexpected error type is returned: %#v.", err)
			}
		})
	}
}
