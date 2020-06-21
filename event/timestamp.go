package event

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// TimeStamp represents slack flavored time representation.
// Slack may pass timestamp in a form of "1355517536.000001," where first preceding integers before dot represents the UNIX timestamp.
// Following integers is used to uniquify the timestamp within a given channel.
// https://api.slack.com/events/message
type TimeStamp struct {
	Time time.Time

	// OriginalValue is exactly the same value slack passes. e.g. "1355517536.000001"
	OriginalValue string
}

// UnmarshalJSON parses a given slack timestamp to time.Time.
// This method is mainly used by encode/json.
func (timeStamp *TimeStamp) UnmarshalJSON(b []byte) error {
	// First accept both string form of "1355517536.000001" and integer form of 1355517536
	var n json.Number
	err := json.Unmarshal(b, &n)
	if err != nil {
		return err
	}
	timeStamp.OriginalValue = n.String()

	// Convert accepted value to time.Time representation.
	i, err := strconv.ParseInt(strings.Split(n.String(), ".")[0], 10, 64)
	if err != nil {
		return err
	}
	timeStamp.Time = time.Unix(i, 0)

	return nil
}

// String returns the original timestamp value given by slack.
func (timeStamp *TimeStamp) String() string {
	return timeStamp.OriginalValue
}

// MarshalText returns the stringified value of slack flavored timestamp.
// To ensure idempotence of marshal/unmarshal, this returns the original value given by slack.
func (timeStamp *TimeStamp) MarshalText() ([]byte, error) {
	return []byte(timeStamp.String()), nil
}
