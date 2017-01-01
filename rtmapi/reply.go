package rtmapi

// WebSocketReply is passed from slack as a reply to client message, and indicates its status.
// https://api.slack.com/rtm#ping_and_pong#handling_responses
type WebSocketReply struct {
	OK        bool       `json:"ok"`
	ReplyTo   uint       `json:"reply_to"`
	TimeStamp *TimeStamp `json:"ts"`
	Text      string     `json:"text"`
}
