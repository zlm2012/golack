package event

// CommonMessage contains some common fields of message event.
// See SubType field to distinguish corresponding event struct.
// https://api.slack.com/events/message#message_subtypes
type CommonMessage struct {
	TypedEvent

	// Regular user message and some miscellaneous message share the common type of "message."
	// So take a look at subtype to distinguish. Regular user message has empty subtype.
	SubType string `json:"subtype"`
}

// MiscMessage represents some minor message events.
// TODO define each one with subtype field. This is just a representation of common subtyped payload
// https://api.slack.com/events/message#message_subtypes
type MiscMessage struct {
	CommonMessage
	TimeStamp *TimeStamp `json:"ts"`
}
