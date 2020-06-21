package eventsapi

import (
	"github.com/oklahomer/golack/event"
	"github.com/oklahomer/golack/testutil"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

var expectedPayloads = map[string]interface{}{
	"url_verification": &URLVerification{
		Type:      "url_verification",
		Challenge: "3eZbrw1aBm2rZgRNFdxV2595E9CY3gmdALWMmHkvFXO7tYXAYM8P",
		Token:     "Jhj5dZrVaK7ZwHHjRyZWjbDl",
	},
	"reaction_added": &EventWrapper{
		outer: &outer{
			Token:       "z26uFbvR1xHJEdHE1OQiO6t8",
			TeamID:      "T061EG9RZ",
			APIAppID:    "A0FFV41KK",
			Type:        "event_callback",
			AuthedUsers: []string{"U061F7AUR"},
			EventID:     "Ev9UQ52YNA",
			EventTime: &event.TimeStamp{
				Time:          time.Unix(1234567890, 0),
				OriginalValue: "1234567890",
			},
		},
		Event: &event.ReactionAdded{
			TypedEvent: event.TypedEvent{
				Type: "reaction_added",
			},
			UserID:   "U061F1EUR",
			Reaction: "slightly_smiling_face",
			Item: &event.Item{
				Type:      "message",
				ChannelID: "C061EG9SL",
				TimeStamp: &event.TimeStamp{
					Time:          time.Unix(1464196127, 0),
					OriginalValue: "1464196127.000002",
				},
			},
			ItemOwnerID: "U0M4RL1NY",
			TimeStamp: &event.TimeStamp{
				Time:          time.Unix(1465244570, 0),
				OriginalValue: "1465244570.336841",
			},
		},
	},
}

func TestDecodePayload(t *testing.T) {
	directory := filepath.Join("..", "testdata", "eventsapi", "decode")
	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Fatalf("Failed to read testdata directory: %s", err.Error())
		}

		if !strings.HasSuffix(path, ".json.golden") {
			// Skip irrelevant file
			return nil
		}

		filename := filepath.Base(path)
		t.Run(filename, func(t *testing.T) {
			input, err := ioutil.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to read file: %s. Error: %s.", path, err.Error())
			}

			eventName := strings.TrimSuffix(filename, ".json.golden")

			request := &SlackRequest{Payload: input}
			payload, err := DecodePayload(request)
			if err != nil {
				t.Fatalf("Failed to decode JSON: %s. Error: %s.", path, err.Error())
			}

			expected, ok := expectedPayloads[eventName]
			if !ok {
				t.Fatalf("Expected payload for %s is not defined in test", eventName)
			}

			if wrapper, ok := payload.(*EventWrapper); ok {
				// Checking request and its []byte payload is a bit troublesome. Skip.
				wrapper.Request = nil
			}
			testutil.Compare([]string{eventName}, reflect.ValueOf(expected), reflect.ValueOf(payload), t)
		})

		return nil
	})
}
