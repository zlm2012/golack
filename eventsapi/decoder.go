package eventsapi

import (
	"encoding/json"
	"fmt"
	"github.com/oklahomer/golack/event"
	"github.com/tidwall/gjson"
	"golang.org/x/xerrors"
)

// https://api.slack.com/events-api#callback_field_overview
type outer struct {
	Token       string           `json:"token"`
	TeamID      string           `json:"team_id"`
	APIAppID    string           `json:"api_app_id"`
	Type        string           `json:"type"`
	AuthedUsers []string         `json:"authed_users"`
	EventID     event.EventID    `json:"event_id"`
	EventTime   *event.TimeStamp `json:"event_time"`
}

// EventWrapper contains given event, metadata and the request.
type EventWrapper struct {
	*outer
	Event   interface{}
	Request *SlackRequest
}

// URLVerification is a special payload for initial configuration.
// When an administrator register an API endpoint to Slack APP configuration page, this payload is sent to the endpoint
// to verify the validity of that endpoint.
//
// This is part of the events list located at https://api.slack.com/events,
// but the structure is defined in this package because this is specifically designed for Events API protocol
// Just like Ping and Pong events are specifically designed for RTM API protocol.
type URLVerification struct {
	Type      string `json:"type"`
	Challenge string `json:"challenge"`
	Token     string `json:"token"`
}

// DecodePayload receives req and decode given event.
// The returning value can be one of *event.URLVerification or *EventWrapper.
// *event.URLVerification can be sent on the initial configuration when an administrator inputs API endpoint to Slack.
func DecodePayload(req *SlackRequest) (interface{}, error) {
	parsed := gjson.ParseBytes(req.Payload)

	typeValue := parsed.Get("type")
	if !typeValue.Exists() {
		return nil, event.NewMalformedPayloadError(fmt.Sprintf("required type field is not given: %s", parsed))
	}

	switch eventType := typeValue.String(); eventType {
	case "url_verification":
		verification := &URLVerification{}
		err := json.Unmarshal(req.Payload, verification)
		if err != nil {
			return nil, xerrors.Errorf("failed to unmarshal JSON: %w", err)
		}
		return verification, nil

	case "event_callback":
		o := &outer{}
		err := json.Unmarshal(req.Payload, o)
		if err != nil {
			return nil, xerrors.Errorf("failed to unmarshal JSON: %w", err)
		}

		// Read the event field that represents the Slack event being sent
		eventValue := parsed.Get("event")
		if !eventValue.Exists() {
			return nil, event.NewMalformedPayloadError(fmt.Sprintf("requred event field is not given: %s", parsed))
		}
		ev, err := event.Map(eventValue)
		if err != nil {
			return nil, err
		}

		// Construct a wrapper object that contains meta, event and request data
		return &EventWrapper{
			outer:   o,
			Event:   ev,
			Request: req,
		}, nil

	default:
		return nil, event.NewUnknownPayloadTypeError(fmt.Sprintf("undefined type of %s is given", eventType))

	}
}
