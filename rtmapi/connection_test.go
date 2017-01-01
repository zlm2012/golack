package rtmapi

import (
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"
)

var webSocketServerAddress string
var once sync.Once

func echoServer(ws *websocket.Conn) {
	defer ws.Close()
	io.Copy(ws, ws)
}

func pingServer(ws *websocket.Conn) {
	defer ws.Close()
	res := &Pong{}
	websocket.JSON.Send(ws, res)
}

func startServer() {
	http.Handle("/echo", websocket.Handler(echoServer))
	http.Handle("/ping", websocket.Handler(pingServer))
	server := httptest.NewServer(nil)
	webSocketServerAddress = server.Listener.Addr().String()
}

func TestConnect(t *testing.T) {
	once.Do(startServer)

	// Establish connection
	url := fmt.Sprintf("ws://%s%s", webSocketServerAddress, "/echo")
	connection, err := Connect(context.TODO(), url)
	if err != nil {
		t.Fatalf("webSocket connection error: %s.", err.Error())
	}

	connWrapper := connection.(*connWrapper)
	if !connWrapper.conn.IsClientConn() {
		t.Fatal("connection is not client originated")
	}
}

func TestConnect_Fail(t *testing.T) {
	once.Do(startServer)

	// Establish connection
	url := fmt.Sprintf("ws://%s%s", webSocketServerAddress, "/undefined_path")
	_, err := Connect(context.TODO(), url)

	if err == nil {
		t.Fatal("expected error is not returned.", err)
	}

	dialErr, ok := err.(*websocket.DialError)
	if !ok {
		t.Fatalf("unexpected error struct is returned: %#v.", err)
	}

	_, ok = dialErr.Err.(*websocket.ProtocolError)
	if !ok {
		t.Fatalf("unexpected error struct is returned: %#v.", err)
	}
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
	once.Do(startServer)

	url := fmt.Sprintf("ws://%s%s", webSocketServerAddress, "/echo")
	conn, err := websocket.Dial(url, "", "http://localhost")
	if err != nil {
		t.Fatal("can't establish connection with test server")
	}
	defer conn.Close()

	channelName := "C12345"
	user := "U6789"
	text := "Hello world!"
	timestamp := 1355517523
	slackTimestamp := fmt.Sprintf("%d.000005", timestamp)
	input := fmt.Sprintf(`{"type": "message", "channel": "%s", "user": "%s", "text": "%s", "ts": "%s"}`, channelName, user, text, slackTimestamp)

	connWrapper := newConnectionWrapper(conn)
	conn.Write([]byte(input))
	decodedPayload, err := connWrapper.Receive()
	if err != nil {
		t.Fatalf("error on payload reception: %s.", err.Error())
	}

	message, ok := decodedPayload.(*Message)
	if !ok {
		t.Fatalf("received payload is not Message: %#v.", message)
	}

	if message.Channel.Name != channelName {
		t.Errorf("expected channel name is not given: %s.", message.Channel.Name)
	}

	if message.Sender != user {
		t.Errorf("expected user is not given: %s.", message.Sender)
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
}

func TestConnWrapper_Send(t *testing.T) {
	once.Do(startServer)

	url := fmt.Sprintf("ws://%s%s", webSocketServerAddress, "/echo")
	conn, err := websocket.Dial(url, "", "http://localhost")
	if err != nil {
		t.Fatal("can't establish connection with test server")
	}
	defer conn.Close()

	connWrapper := newConnectionWrapper(conn)
	if err := connWrapper.Send(&Channel{Name: "dummy channel"}, "hello"); err != nil {
		t.Errorf("error on sending message over WebSocket connection. %#v.", err)
	}
}

func TestConnWrapper_Ping(t *testing.T) {
	once.Do(startServer)

	url := fmt.Sprintf("ws://%s%s", webSocketServerAddress, "/ping")
	conn, err := websocket.Dial(url, "", "http://localhost")
	if err != nil {
		t.Fatal("can't establish connection with test server")
	}
	defer conn.Close()

	connWrapper := newConnectionWrapper(conn)
	if err := connWrapper.Ping(); err != nil {
		t.Errorf("error on sending message over WebSocket connection. %#v.", err)
	}
}

func TestConnWrapper_Close(t *testing.T) {
	once.Do(startServer)

	url := fmt.Sprintf("ws://%s%s", webSocketServerAddress, "/echo")
	conn, err := websocket.Dial(url, "", "http://localhost")
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
}

func TestDecodePayload(t *testing.T) {
	type output struct {
		payload reflect.Type
		err     interface{}
	}
	var decodeTests = []struct {
		input  string
		output output
	}{
		{
			`{"type": "message", "channel": "C2147483705", "user": "U2147483697", "text": "Hello, world!", "ts": "1355517523.000005", "edited": { "user": "U2147483697", "ts": "1355517536.000001"}}`,
			output{
				reflect.TypeOf(&Message{}),
				nil,
			},
		},
		{
			// type is valid and hence mapped to Message, but can not be parsed since the timestamp format is illegal.
			`{"type": "message", "ts": "invalid timestamp"}`,
			output{
				nil,
				reflect.TypeOf(&MalformedPayloadError{}),
			},
		},
		{
			`{"type": "message", "subtype": "channel_join", "text": "<@UXXXXX|bobby> has joined the channel", "ts": "1403051575.000407", "user": "U023BECGF"}`,
			output{
				reflect.TypeOf(&MiscMessage{}),
				nil,
			},
		},
		{
			// invalid type
			`{"type": "unsupportedEventType"}`,
			output{
				nil,
				reflect.TypeOf(&MalformedPayloadError{}),
			},
		},
		{
			`{"ok": true, "reply_to": 1, "ts": "1355517523.000005", "text": "Hello world"}`,
			output{
				reflect.TypeOf(&WebSocketReply{}),
				nil,
			},
		},
		{
			// required fields are not given
			`{"what": true}`,
			output{
				nil,
				reflect.TypeOf(&MalformedPayloadError{}),
			},
		},
		{
			// not valid json structure
			`malformedJson`,
			output{
				nil,
				reflect.TypeOf(&MalformedPayloadError{}),
			},
		},
		{
			"　",
			output{
				nil,
				ErrEmptyPayload,
			},
		},
		{
			"\r",
			output{
				nil,
				ErrEmptyPayload,
			},
		},
	}

	for i, testSet := range decodeTests {
		testCnt := i + 1
		inputByte := []byte(testSet.input)
		payload, err := decodePayload(inputByte)

		if testSet.output.payload != reflect.TypeOf(payload) {
			t.Errorf("Test No. %d. expected return type of %s, but was %#v", testCnt, testSet.output.payload.Name(), err)
		}
		if e := testSet.output.err; e != nil {
			if reflect.TypeOf(e) == reflect.TypeOf(errors.New("Dummy")) {
				// pre-declared error instance is returned.
				if e != err {
					t.Errorf("expected error is not returned. test #%d. error: %#v.", testCnt, err)
				}
			} else {
				// new error instance of specific error struct is returned.
				if e != reflect.TypeOf(err) {
					t.Errorf("unexpected error type is returned on test #%d. error: %#v.", testCnt, err)
				}
			}
		}
	}
}