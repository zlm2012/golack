package eventsapi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	SlackSignatureHeaderName        = "X-Slack-Signature"
	SlackRequestTimestampHeaderName = "X-Slack-Request-Timestamp"
)

// SlackRequest represents the request being sent from Slack.
// This not only includes the payload given as request body, but also includes a group of header values Slack defines.
// See https://api.slack.com/events-api for detailed protocol.
type SlackRequest struct {
	Signature string
	TimeStamp time.Time
	Payload   []byte
}

// NewSlackRequest receives r and instantiate SlackRequest.
func NewSlackRequest(r *http.Request) (*SlackRequest, error) {
	defer r.Body.Close()

	signature := r.Header.Get(SlackSignatureHeaderName)
	timestamp := r.Header.Get(SlackRequestTimestampHeaderName)

	if signature == "" {
		return nil, &BadRequestError{Err: fmt.Sprintf("required %s header is absent", SlackSignatureHeaderName)}
	}

	if timestamp == "" {
		return nil, &BadRequestError{Err: fmt.Sprintf("required %s header is absent", SlackRequestTimestampHeaderName)}
	}

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read payload: %w", err)
	}

	ts, err := strconv.Atoi(timestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timestamp: %w", err)
	}

	return &SlackRequest{
		Signature: signature,
		TimeStamp: time.Unix(int64(ts), 0),
		Payload:   payload,
	}, nil
}

// BadRequestError implies the given request is invalid.
type BadRequestError struct {
	Err string
}

// Error returns detailed error state.
func (e *BadRequestError) Error() string {
	return e.Err
}
