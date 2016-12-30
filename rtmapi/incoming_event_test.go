package rtmapi

import (
	"reflect"
	"testing"
)

func TestDecodeEvent(t *testing.T) {
	type output struct {
		payload reflect.Type
		err     error
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
			`{"type": "message", "subtype": "channel_join", "text": "<@UXXXXX|bobby> has joined the channel", "ts": "1403051575.000407", "user": "U023BECGF"}`,
			output{
				reflect.TypeOf(&MiscMessage{}),
				nil,
			},
		},
		{
			`{"type": "", "channel": "C2147483705"}`,
			output{
				nil,
				ErrUnsupportedEventType,
			},
		},
		{
			`{"type": "foo", "channel": "C2147483705"}`,
			output{
				nil,
				ErrUnsupportedEventType,
			},
		},
		{
			`{"channel": "C2147483705"}`,
			output{
				nil,
				ErrEventTypeNotGiven,
			},
		},
	}

	for i, testSet := range decodeTests {
		testCnt := i + 1
		inputByte := []byte(testSet.input)
		event, err := DecodeEvent(inputByte)

		if testSet.output.payload != reflect.TypeOf(event) {
			t.Errorf("Test No. %d. expected return type of %s, but was %#v", testCnt, testSet.output.payload.Name(), err)
		}
		if testSet.output.err != err {
			t.Errorf("Test No. %d. Expected return error of %#v, but was %#v", testCnt, testSet.output.err, err)
		}
	}
}
