package testutil

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"net/http/httptest"
)

type Pong struct {
	Type    string `json:"type"`
	ReplyTo uint   `json:"reply_to"`
}

func echoServer(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(fmt.Errorf("failed to upgrade protocol: %s", err.Error()))
	}
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			_, ok := err.(*websocket.CloseError)
			if !ok {
				panic(fmt.Errorf("failed to receive message: %s", err.Error()))
			}
			break
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			panic(fmt.Errorf("failed to send message. type: %d. error: %s", mt, err.Error()))
		}
	}
}

func pingServer(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(fmt.Errorf("failed to upgrade protocol: %s", err.Error()))
	}
	defer c.Close()

	res := &Pong{}
	c.WriteJSON(res)
}

func RunWithWebSocket(fnc func(addr net.Addr)) {
	// Setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/echo", echoServer)
	mux.HandleFunc("/ping", pingServer)
	server := httptest.NewServer(mux)

	// Close after test
	defer server.Close()

	// Run test with WebSocket server
	fnc(server.Listener.Addr())
}
