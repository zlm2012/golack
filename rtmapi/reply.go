package rtmapi

// WebSocketReply represents a response from Slack.
// No matter if the response represents successful operation or not,
// the reply always contains this struct.
//
// When the response indicates successful state, the payload should be mapped to WebSocketOKReply;
// while the payload should be mapped to WebSocketNGReply on failure state.
// https://api.slack.com/rtm#handling_responses
type WebSocketReply struct {
	OK      bool `json:"ok"`
	ReplyTo uint `json:"reply_to"`
}

// WebSocketOKReply represents a successful operation.
type WebSocketOKReply struct {
	WebSocketReply
	TimeStamp *TimeStamp `json:"ts"`
	Text      string     `json:"text"`
}

// WebSocketOKReply represents a failed operation.
type WebSocketNGReply struct {
	WebSocketReply
	ErrorReason struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
	} `json:"error"`
}
