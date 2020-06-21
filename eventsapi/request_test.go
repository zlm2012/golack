package eventsapi

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestNewSlackRequest(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		signature := "signature"
		now := time.Unix(1234567890, 0)
		payload := "payload"
		req := &http.Request{
			Header: map[string][]string{
				SlackSignatureHeaderName:        {signature},
				SlackRequestTimestampHeaderName: {strconv.FormatInt(now.Unix(), 10)},
			},
			Body: ioutil.NopCloser(strings.NewReader(payload)),
		}

		request, err := NewSlackRequest(req)

		if err != nil {
			t.Fatalf("Unexpected error is returned: %s.", err.Error())
		}

		if request.Signature != signature {
			t.Errorf("Expected signature value of %s, but was %s.", signature, request.Signature)
		}

		if !request.TimeStamp.Equal(now) {
			t.Errorf("Expected timestamp to be %s, but was %s.", now, request.TimeStamp)
		}

		if string(request.Payload) != payload {
			t.Errorf("Expected payload to be %s, but was %s.", payload, request.Payload)
		}
	})

	t.Run("missing signature", func(t *testing.T) {
		now := time.Unix(1234567890, 0)
		payload := "payload"
		req := &http.Request{
			Header: map[string][]string{
				SlackRequestTimestampHeaderName: {strconv.FormatInt(now.Unix(), 10)},
			},
			Body: ioutil.NopCloser(strings.NewReader(payload)),
		}

		_, err := NewSlackRequest(req)

		if err == nil {
			t.Fatal("Expected error is not returned.")
		}

		typedError, ok := err.(*BadRequestError)
		if !ok {
			t.Errorf("Expected *BadRequestError, but was %v", err)
		}

		if !strings.Contains(typedError.Error(), SlackSignatureHeaderName) {
			t.Errorf("Returned error message does not contain header name")
		}
	})

	t.Run("missing time", func(t *testing.T) {
		signature := "signature"
		payload := "payload"
		req := &http.Request{
			Header: map[string][]string{
				SlackSignatureHeaderName: {signature},
			},
			Body: ioutil.NopCloser(strings.NewReader(payload)),
		}

		_, err := NewSlackRequest(req)

		if err == nil {
			t.Fatal("Expected error is not returned.")
		}

		typedError, ok := err.(*BadRequestError)
		if !ok {
			t.Errorf("Expected *BadRequestError, but was %v", err)
		}

		if !strings.Contains(typedError.Error(), SlackRequestTimestampHeaderName) {
			t.Error("Returned error message does not contain header name")
		}
	})

	t.Run("invalid time", func(t *testing.T) {
		signature := "signature"
		payload := "payload"
		req := &http.Request{
			Header: map[string][]string{
				SlackSignatureHeaderName:        {signature},
				SlackRequestTimestampHeaderName: {"invalid"},
			},
			Body: ioutil.NopCloser(strings.NewReader(payload)),
		}

		_, err := NewSlackRequest(req)

		if err == nil {
			t.Fatal("Expected error is not returned.")
		}
	})
}
