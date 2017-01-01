package rtmapi

// Hello event is sent from slack when WebSocket connection is successfully established.
// https://api.slack.com/events/hello
type Hello struct {
	CommonEvent
}

// TeamMigrationStarted is sent when chat group is migrated between servers.
// "The WebSocket connection will close immediately after it is sent.
// *snip* By the time a client has reconnected the process is usually complete, so the impact is minimal."
// https://api.slack.com/events/team_migration_started
type TeamMigrationStarted struct {
	CommonEvent
}

// Pong is given when client send Ping.
// https://api.slack.com/rtm#ping_and_pong
type Pong struct {
	CommonEvent
	ReplyTo uint `json:"reply_to"`
}

// IncomingChannelEvent represents any event occurred in a specific channel.
// This can be a part of other event such as message.
type IncomingChannelEvent struct {
	CommonEvent
	Channel *Channel `json:"channel"`
}

// Message represent message event on RTM.
// https://api.slack.com/events/message
// This implements Input interface.
//  {
//      "type": "message",
//      "channel": "C2147483705",
//      "user": "U2147483697",
//      "text": "Hello, world!",
//      "ts": "1355517523.000005",
//      "edited": {
//          "user": "U2147483697",
//          "ts": "1355517536.000001"
//      }
//  }
type Message struct {
	IncomingChannelEvent
	Sender    string     `json:"user"`
	Text      string     `json:"text"`
	TimeStamp *TimeStamp `json:"ts"`
}

// MiscMessage represents some minor message events.
// TODO define each one with subtype field. This is just a representation of common subtyped payload
// https://api.slack.com/events/message#message_subtypes
type MiscMessage struct {
	CommonMessage
	TimeStamp *TimeStamp `json:"ts"`
}
