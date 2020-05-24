package rtmapi

import "github.com/oklahomer/golack/event"

// Reply represents a response from Slack.
// No matter if the response represents successful operation or not,
// the reply always contains this struct.
//
// When the response indicates successful state, the payload should be mapped to OKReply;
// while the payload should be mapped to WebSocketNGReply on failure state.
// https://api.slack.com/rtm#handling_responses
type Reply struct {
	OK      bool `json:"ok"`
	ReplyTo uint `json:"reply_to"`
}

// OKReply represents a successful operation.
type OKReply struct {
	Reply
	TimeStamp *event.TimeStamp `json:"ts"`
	Text      string           `json:"text"`
}

// NGReply represents a failed operation.
type NGReply struct {
	Reply
	Error ReplyErrorReason `json:"error"`
}

type ReplyErrorReason struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}
