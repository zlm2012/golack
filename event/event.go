package event

import (
	"encoding/json"
	"golang.org/x/xerrors"
)

// TypedEvent takes care of events that have "type" field in its JSON representation.
//
// RTM API document, https://api.slack.com/rtm#events, states as follows:
// "Every event has a type property which describes the type of event."
//
// Similarly, Events API document, https://api.slack.com/events-api#event_type_structure, states as follows:
// "The specific name of the event described by its adjacent fields. This field is included with every inner event type."
type TypedEvent struct {
	Type string `json:"type"`
}

var _ Typer = (*TypedEvent)(nil)

// EventType returns the type of event the payload is representing.
func (t TypedEvent) EventType() string {
	return t.Type
}

// Typer defines an interface that all typed event must satisfy.
type Typer interface {
	EventType() string
}

// https://api.slack.com/events
// AccountsChanged event indicates that the list of accounts a user is signed into has changed.
// https://api.slack.com/events/accounts_changed
type AccountsChanged struct {
	TypedEvent
}

// AppHomeOpened event represents an event a user clicked into App Home
// https://api.slack.com/events/app_home_opened
//
// Block object:https://api.slack.com/reference/block-kit/blocks#input
// View object: https://api.slack.com/reference/interaction-payloads/views
type AppHomeOpened struct {
	TypedEvent
	UserID         UserID     `json:"user"`
	ChannelID      ChannelID  `json:"channel"`
	EventTimeStamp *TimeStamp `json:"event_ts"`
	Tab            string     `json:"tab"`
	View           *View      `json:"view"`
}

// AppMention allows your app to subscribe to message events that directly mention your bot user.
// https://api.slack.com/events/app_mention
type AppMention struct {
	TypedEvent
	UserID         UserID     `json:"user"`
	Text           string     `json:"text"`
	TimeStamp      *TimeStamp `json:"ts"`
	ChannelID      ChannelID  `json:"channel"`
	EventTimeStamp *TimeStamp `json:"event_ts"`
}

// AppRateLimited is only dispatched when your app is rate limited on the Events API.
// https://api.slack.com/events/app_rate_limited
type AppRateLimited struct {
	TypedEvent
	Token             string     `json:"token"`
	TeamID            string     `json:"team_id"`
	MinuteRateLimited *TimeStamp `json:"minute_rate_limited"`
	APIAppID          AppID      `json:"api_app_id"`
}

// AppRequested contains information about an app that a user on a team has requested to install.
// https://api.slack.com/events/app_requested
type AppRequested struct {
	TypedEvent
	AppRequest *AppRequest `json:"app_request"`
}

// AppUninstalled is sent via subscription whenever a Slack app is completely uninstalled.
// https://api.slack.com/events/app_uninstalled
type AppUninstalled struct {
	TypedEvent
}

type BotAdded struct {
	TypedEvent
	Bot *Bot `json:"bot"`
}

type BotChanged struct {
	TypedEvent
	Bot *Bot `json:"bot"`
}

// CallRejected is sent if the user rejects the Call.
// https://api.slack.com/events/call_rejected
type CallRejected struct {
	TypedEvent
	CallID           CallID    `json:"call_id"`
	UserID           UserID    `json:"user_id"`
	ChannelID        ChannelID `json:"channel_id"`
	ExternalUniqueID string    `json:"external_unique_id"`
}

type ChannelArchived struct {
	TypedEvent
	ChannelID ChannelID `json:"channel"`
	UserID    UserID    `json:"user"`
}

type ChannelCreated struct {
	TypedEvent
	Channel *struct {
		ID        ChannelID  `json:"id"`
		Name      string     `json:"name"`
		Created   *TimeStamp `json:"created"`
		CreatorID UserID     `json:"creator"`
	} `json:"channel"`
}

type ChannelDeleted struct {
	TypedEvent
	ChannelID ChannelID `json:"channel"`
}

type ChannelHistoryChanged struct {
	TypedEvent
	ChangedHistory
}

type ChannelJoined struct {
	TypedEvent
	ChannelID ChannelID `json:"channel"`
}

type ChannelLeft struct {
	TypedEvent
	ChannelID ChannelID `json:"channel"`
}

type ChannelMarked struct {
	TypedEvent
	MarkedAsRead
}

type ChannelRenamed struct {
	TypedEvent
	Channel *struct {
		ID      ChannelID  `json:"id"`
		Name    string     `json:"name"`
		Created *TimeStamp `json:"created"`
	} `json:"channel"`
}

// ChannelShared is sent to all event subscriptions when a new shared channel is created or a channel is converted into a shared channel.
// https://api.slack.com/events/channel_shared
type ChannelShared struct {
	TypedEvent
	ConnectedTeamID TeamID     `json:"connected_team_id"`
	ChannelID       ChannelID  `json:"channel"`
	EventTimeStamp  *TimeStamp `json:"event_ts"`
}

type ChannelUnarchived struct {
	TypedEvent
	ChannelID ChannelID `json:"channel"`
	UserID    UserID    `json:"user"`
}

// ChannelUnshared is sent to all event subscriptions when an external workspace has been removed from an existing shared channel.
// https://api.slack.com/events/channel_unshared
type ChannelUnshared struct {
	TypedEvent
	PreviouslyConnectedTeamID TeamID     `json:"previously_connected_team_id"`
	ChannelID                 ChannelID  `json:"channel"`
	IsExtShared               bool       `json:"is_ext_shared"`
	EventTimeStamp            *TimeStamp `json:"event_ts"`
}

type CommandsChanged struct {
	TypedEvent
	TimeStamp *TimeStamp `json:"event_ts"`
}

type DNDUpdated struct {
	TypedEvent
	UserID    UserID `json:"user"`
	DNDStatus *struct {
		Enabled            bool       `json:"dnd_enabled"`
		NextStartTimeStamp *TimeStamp `json:"next_dnd_start_ts"`
		NextEndTimeStamp   *TimeStamp `json:"next_dnd_end_ts"`
		SnoozeEnabled      bool       `json:"snooze_enabled"`
		SnoozeEndTimeStamp *TimeStamp `json:"snooze_endtime"`
	} `json:"dnd_status"`
}

type DNDUpdatedUser struct {
	TypedEvent
	UserID    UserID `json:"user"`
	DNDStatus *struct {
		Enabled            bool       `json:"dnd_enabled"`
		NextStartTimeStamp *TimeStamp `json:"next_dnd_start_ts"`
		NextEndTimeStamp   *TimeStamp `json:"next_dnd_end_ts"`
	} `json:"dnd_status"`
}

type EmailDomainChanged struct {
	TypedEvent
	EmailDomain string     `json:"email_domain"`
	TimeStamp   *TimeStamp `json:"event_ts"`
}

type EmojiChanged struct {
	TypedEvent
	Subtype   string     `json:"subtype"` // TODO add/remove
	Names     []string   `json:"names"`
	TimeStamp *TimeStamp `json:"event_ts"`
}

// ExternalOrgMigrationFinished is sent to all connections when an external workspace completes to migrate to an Enterprise Grid.
// https://api.slack.com/events/external_org_migration_finished
type ExternalOrgMigrationFinished struct {
	TypedEvent
	Team *struct {
		ID          TeamID `json:"id"`
		IsMigrating bool   `json:"is_migrating"`
	} `json:"team"`
	DateStarted  *TimeStamp `json:"date_started"`
	DateFinished *TimeStamp `json:"date_finished"`
}

// ExternalOrgMigrationStarted is sent to all connections when an external workspace begins to migrate to an Enterprise Grid.
// https://api.slack.com/events/external_org_migration_started
type ExternalOrgMigrationStarted struct {
	TypedEvent
	Team *struct {
		ID          TeamID `json:"id"`
		IsMigrating bool   `json:"is_migrating"`
	} `json:"team"`
	DateStarted *TimeStamp `json:"date_started"`
}

type FileChanged struct {
	TypedEvent
	FileID FileID `json:"file_id"`
	File   *struct {
		// https://api.slack.com/events/file_change
		// This value is identical to FileID field
		ID FileID `json:"id"`
	} `json:"file"`
}

type FileCommentAdded struct {
	TypedEvent
	FileID FileID `json:"file_id"`
	File   *struct {
		// https://api.slack.com/events/file_change
		// This value is identical to FileID field
		ID FileID `json:"id"`
	} `json:"file"`
	Comment *Comment `json:"comment"`
}

type FileCommentDeleted struct {
	TypedEvent
	FileID FileID `json:"file_id"`
	File   *struct {
		// https://api.slack.com/events/file_change
		// This value is identical to FileID field
		ID FileID `json:"id"`
	} `json:"file"`
	CommentID string `json:"comment"`
}

type FileCommentEdited struct {
	TypedEvent
	FileID FileID `json:"file_id"`
	File   *struct {
		// https://api.slack.com/events/file_change
		// This value is identical to FileID field
		ID FileID `json:"id"`
	} `json:"file"`
	Comment *Comment `json:"comment"`
}

type FileCreated struct {
	TypedEvent
	FileID FileID `json:"file_id"`
	File   *struct {
		// https://api.slack.com/events/file_change
		// This value is identical to FileID field
		ID FileID `json:"id"`
	} `json:"file"`
}

type FileDeleted struct {
	TypedEvent
	FileID    FileID     `json:"file_id"`
	TimeStamp *TimeStamp `json:"event_ts"`
}

type FilePublished struct {
	TypedEvent
	FileID FileID `json:"file_id"`
	File   *struct {
		// https://api.slack.com/events/file_change
		// This value is identical to FileID field
		ID FileID `json:"id"`
	} `json:"file"`
}

type FileShared struct {
	TypedEvent
	FileID FileID `json:"file_id"`
	File   *struct {
		// https://api.slack.com/events/file_change
		// This value is identical to FileID field
		ID FileID `json:"id"`
	} `json:"file"`
}

type FileUnshared struct {
	TypedEvent
	FileID FileID `json:"file_id"`
	File   *struct {
		// https://api.slack.com/events/file_change
		// This value is identical to FileID field
		ID FileID `json:"id"`
	} `json:"file"`
}

type GoodBye struct {
	TypedEvent
}

// GridMigrationFinished is sent via subscription whenever your app is installed by completes migration to Enterprise Grid.
// https://api.slack.com/events/grid_migration_finished
type GridMigrationFinished struct {
	TypedEvent
	EnterpriseID string `json:"enterprise_id"`
}

// GridMigrationStarted is sent via subscription whenever your app is installed by begins to migrate to an Enterprise Grid.
// https://api.slack.com/events/grid_migration_started
type GridMigrationStarted struct {
	TypedEvent
	EnterpriseID string `json:"enterprise_id"`
}

type GroupArchived struct {
	TypedEvent
	ChannelID ChannelID `json:"channel"`
}

type GroupClosed struct {
	TypedEvent
	UserID    UserID    `json:"user"`
	ChannelID ChannelID `json:"channel"`
}

type GroupDeleted struct {
	TypedEvent
	ChannelID ChannelID `json:"channel"`
}

type GroupHistoryChanged struct {
	TypedEvent
	ChangedHistory
}

type GroupJoined struct {
	TypedEvent
	ChannelID ChannelID `json:"channel"`
}

type GroupLeft struct {
	TypedEvent
	ChannelID ChannelID `json:"channel"`
}

type GroupMarked struct {
	TypedEvent
	MarkedAsRead
}

type GroupOpened struct {
	TypedEvent
	UserID    UserID    `json:"user"`
	ChannelID ChannelID `json:"channel"`
}

type GroupRenamed struct {
	TypedEvent
	Channel *struct {
		ID      ChannelID  `json:"id"`
		Name    string     `json:"name"`
		Created *TimeStamp `json:"created"`
	} `json:"channel"`
}

type GroupUnarchived struct {
	TypedEvent
	ChannelID ChannelID `json:"channel"`
}

// Hello event is sent from slack when WebSocket connection is successfully established.
// https://api.slack.com/events/hello
type Hello struct {
	TypedEvent
}

type IMClosed struct {
	TypedEvent
	UserID    UserID    `json:"user"`
	ChannelID ChannelID `json:"channel"`
}

type IMCreated struct {
	TypedEvent
	UserID  UserID `json:"user"`
	Channel *struct {
		ID ChannelID `json:"id"`
	} `json:"channel"`
}

type IMHistoryChanged struct {
	TypedEvent
	ChangedHistory
}

type IMMarked struct {
	TypedEvent
	MarkedAsRead
}

type IMOpened struct {
	TypedEvent
	UserID    UserID    `json:"user"`
	ChannelID ChannelID `json:"channel"`
}

// InviteRequested is sent when a user request an invite.
// https://api.slack.com/events/invite_requested
type InviteRequested struct {
	TypedEvent
	InviteRequest *struct {
		ID            string      `json:"id"`
		Email         string      `json:"email"`
		DateCreated   *TimeStamp  `json:"date_created"`
		RequesterIDs  []UserID    `json:"requester_ids"`
		ChannelIDs    []ChannelID `json:"channel_ids"`
		InviteType    string      `json:"invite_type"`
		RealName      string      `json:"real_name"`
		DateExpire    *TimeStamp  `json:"date_expire"`
		RequestReason string      `json:"request_reason"`
		Team          *struct {
			ID     TeamID `json:"id"`
			Name   string `json:"name"`
			Domain string `json:"domain"`
		} `json:"team"`
	} `json:"invite_request"`
}

// LinkShared is sent to track a specific URL domain.
// https://api.slack.com/events/link_shared
type LinkShared struct {
	TypedEvent
	ChannelID        ChannelID  `json:"channel"`
	UserID           UserID     `json:"user"`
	MessageTimeStamp *TimeStamp `json:"message_ts"`
	ThreadTimeStamp  *TimeStamp `json:"thread_ts"`
	Links            []*struct {
		Domain string `json:"domain"`
		URL    string `json:"url"`
	} `json:"links"`
}

type PresenceManuallyChanged struct {
	TypedEvent
	Presence string `json:"presence"` // TODO Actual values other than "away"
}

type MemberJoinedChannel struct {
	TypedEvent
	UserID      UserID    `json:"user"`
	ChannelID   ChannelID `json:"channel"`
	ChannelType string    `json:"channel_type"` // C or G. ref. https://api.slack.com/events/member_joined_channel
	TeamID      TeamID    `json:"team"`
	InviterID   UserID    `json:"inviter"` // Empty when the user joins by herself. ref. https://api.slack.com/events/member_joined_channel
}

// MemberLeftChannel is sent to all websocket connections and event subscriptions when users leave public or private channels.
// https://api.slack.com/events/member_left_channel
type MemberLeftChannel struct {
	TypedEvent
	UserID      UserID    `json:"user"`
	ChannelID   ChannelID `json:"channel"`
	ChannelType string    `json:"channel_type"` // TODO C or G https://api.slack.com/events/member_left_channel
	TeamID      TeamID    `json:"team"`
}

// Message represents a message event on RTM.
// https://api.slack.com/events/message
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
	TypedEvent
	ChannelID       ChannelID  `json:"channel"`
	UserID          UserID     `json:"user"`
	Text            string     `json:"text"`
	TimeStamp       *TimeStamp `json:"ts"`
	ThreadTimeStamp *TimeStamp `json:"thread_ts"` // https://api.slack.com/docs/message-threading
	Edited          *struct {
		UserID    UserID     `json:"user"`
		TimeStamp *TimeStamp `json:"ts"`
	} `json:"edited"`
	Replies []*struct {
		UserID    UserID     `json:"user"`
		TimeStamp *TimeStamp `json:"ts"`
	} `json:"replies"`
}

// ChannelTypeMessage represents a message event on Events API.
//
// See below documents:
//   - https://api.slack.com/events/message
//   - https://api.slack.com/events/message.app_home
//   - https://api.slack.com/events/message.channels
//   - https://api.slack.com/events/message.groups
//   - https://api.slack.com/events/message.im
//   - https://api.slack.com/events/message.mpim
type ChannelMessage struct {
	TypedEvent
	ChannelID       ChannelID  `json:"channel"`
	UserID          UserID     `json:"user"`
	Text            string     `json:"text"`
	TimeStamp       *TimeStamp `json:"ts"`
	ThreadTimeStamp *TimeStamp `json:"thread_ts"` // https://api.slack.com/docs/message-threading
	Edited          *struct {
		UserID    UserID     `json:"user"`
		TimeStamp *TimeStamp `json:"ts"`
	} `json:"edited"`
	Replies []*struct {
		UserID    UserID     `json:"user"`
		TimeStamp *TimeStamp `json:"ts"`
	} `json:"replies"`
	EventTimeStamp *TimeStamp `json:"event_ts"`
	ChannelType    string     `json:"channel_type"`
}

type MessageBotMessage struct {
	TypedEvent
	SubType   string     `json:"subtype"`
	TimeStamp *TimeStamp `json:"ts"`
	Text      string     `json:"text"`
	BotID     BotID      `json:"bot_id"`
	UserName  string     `json:"username"`
	Icon      *BotIcon   `json:"icons"`
}

type MessageChannelArchive struct {
	TypedEvent
	SubType   string     `json:"subtype"`
	TimeStamp *TimeStamp `json:"ts"`
	Text      string     `json:"text"`
	UserID    UserID     `json:"user"`
}

type MessageChannelJoin struct {
	TypedEvent
	SubType   string     `json:"subtype"`
	TimeStamp *TimeStamp `json:"ts"`
	Text      string     `json:"text"`
	UserID    UserID     `json:"user"`
	InviterID UserID     `json:"inviter"`
}

type MessageChannelLeave struct {
	TypedEvent
	SubType   string     `json:"subtype"`
	TimeStamp *TimeStamp `json:"ts"`
	Text      string     `json:"text"`
	UserID    UserID     `json:"user"`
}

type MessageChannelName struct {
	TypedEvent
	SubType   string     `json:"subtype"`
	TimeStamp *TimeStamp `json:"ts"`
	Text      string     `json:"text"`
	UserID    UserID     `json:"user"`
	OldName   string     `json:"old_name"`
	Name      string     `json:"name"`
}

type MessageChannelPurpose struct {
	TypedEvent
	SubType   string     `json:"subtype"`
	TimeStamp *TimeStamp `json:"ts"`
	Text      string     `json:"text"`
	UserID    UserID     `json:"user"`
	Purpose   string     `json:"purpose"`
}

type MessageChannelTopic struct {
	TypedEvent
	SubType   string     `json:"subtype"`
	TimeStamp *TimeStamp `json:"ts"`
	Text      string     `json:"text"`
	UserID    UserID     `json:"user"`
	Topic     string     `json:"topic"`
}

type MessageChannelUnarchive struct {
	TypedEvent
	SubType   string     `json:"subtype"`
	TimeStamp *TimeStamp `json:"ts"`
	Text      string     `json:"text"`
	UserID    UserID     `json:"user"`
}

type MessageEKMAccessDenied struct {
	TypedEvent
	SubType   string     `json:"subtype"`
	TimeStamp *TimeStamp `json:"ts"`
	Text      string     `json:"text"`
	UserID    UserID     `json:"user"`
}

type MessageFileComment struct {
	TypedEvent
	SubType   string     `json:"subtype"`
	TimeStamp *TimeStamp `json:"ts"`
	Text      string     `json:"text"`
	File      *File      `json:"file"`
	Comment   *Comment   `json:"comment"`
}

type MessageFileMention struct {
	TypedEvent
	SubType   string     `json:"subtype"`
	TimeStamp *TimeStamp `json:"ts"`
	Text      string     `json:"text"`
	File      *File      `json:"file"`
	UserID    UserID     `json:"user"`
}

type MessageFileShare struct {
	TypedEvent
	SubType        string     `json:"subtype"`
	Files          []*File    `json:"files"`
	TimeStamp      *TimeStamp `json:"ts"`
	Text           string     `json:"text"`
	File           *File      `json:"file"`
	UserID         UserID     `json:"user"`
	Upload         bool       `json:"upload"`
	DisplayAsBot   bool       `json:"display_as_bot"`
	BotID          BotID      `json:"bot_id"`
	EventTimeStamp *TimeStamp `json:"event_ts"`
	ChannelType    string     `json:"channel_type"`
}

type MessageGroupArchive struct {
	TypedEvent
	SubType   string     `json:"subtype"`
	TimeStamp *TimeStamp `json:"ts"`
	Text      string     `json:"text"`
	UserID    UserID     `json:"user"`
	MemberIDs []UserID   `json:"members"`
}

type MessageGroupJoin struct {
	TypedEvent
	SubType   string     `json:"subtype"`
	TimeStamp *TimeStamp `json:"ts"`
	Text      string     `json:"text"`
	UserID    UserID     `json:"user"`
}

type MessageGroupLeave struct {
	TypedEvent
	SubType   string     `json:"subtype"`
	TimeStamp *TimeStamp `json:"ts"`
	Text      string     `json:"text"`
	UserID    UserID     `json:"user"`
}

type MessageGroupName struct {
	TypedEvent
	SubType   string     `json:"subtype"`
	TimeStamp *TimeStamp `json:"ts"`
	Text      string     `json:"text"`
	UserID    UserID     `json:"user"`
	OldName   string     `json:"old_name"`
	Name      string     `json:"name"`
}

type MessageGroupPurpose struct {
	TypedEvent
	SubType   string     `json:"subtype"`
	TimeStamp *TimeStamp `json:"ts"`
	Text      string     `json:"text"`
	UserID    UserID     `json:"user"`
	Purpose   string     `json:"purpose"`
}

type MessageGroupTopic struct {
	TypedEvent
	SubType   string     `json:"subtype"`
	TimeStamp *TimeStamp `json:"ts"`
	Text      string     `json:"text"`
	UserID    UserID     `json:"user"`
	Topic     string     `json:"topic"`
}

type MessageGroupUnarchive struct {
	TypedEvent
	SubType   string     `json:"subtype"`
	TimeStamp *TimeStamp `json:"ts"`
	Text      string     `json:"text"`
	UserID    UserID     `json:"user"`
}

type MessageMeMessage struct {
	TypedEvent
	SubType   string     `json:"subtype"`
	TimeStamp *TimeStamp `json:"ts"`
	Text      string     `json:"text"`
	UserID    UserID     `json:"user"`
	ChannelID ChannelID  `json:"channel"`
}

type MessageMessageChanged struct {
	TypedEvent
	SubType   string     `json:"subtype"`
	TimeStamp *TimeStamp `json:"ts"`
	ChannelID ChannelID  `json:"channel"`
	Message   *Message   `json:"message"`
	Hidden    bool       `json:"hidden"`
}

type MessageMessageDeleted struct {
	TypedEvent
	SubType          string     `json:"subtype"`
	TimeStamp        *TimeStamp `json:"ts"`
	ChannelID        ChannelID  `json:"channel"`
	Hidden           bool       `json:"hidden"`
	DeletedTimeStamp *TimeStamp `json:"deleted_ts"`
}

type MessageMessageReplied struct {
	TypedEvent
	SubType        string     `json:"subtype"`
	TimeStamp      *TimeStamp `json:"ts"`
	ChannelID      ChannelID  `json:"channel"`
	Hidden         bool       `json:"hidden"`
	EventTimeStamp *TimeStamp `json:"event_ts"`
	Message        *Message   `json:"message"`
}

type MessagePinnedItem struct {
	TypedEvent
	SubType   string     `json:"subtype"`
	TimeStamp *TimeStamp `json:"ts"`
	UserID    UserID     `json:"user"`
	Text      string     `json:"text"`
	ChannelID ChannelID  `json:"channel"`
	ItemType  ItemType   `json:"item_type"`
	Item      *Item      `json:"item"`
}

type MessageThreadBroadcast struct {
	TypedEvent
	SubType         string     `json:"subtype"`
	TimeStamp       *TimeStamp `json:"ts"`
	UserID          UserID     `json:"user"`
	Text            string     `json:"text"`
	Root            *Message   `json:"root"`
	ThreadTimeStamp *TimeStamp `json:"thread_ts"`
}

type MessageUnpinnedItem struct {
	TypedEvent
	SubType   string     `json:"subtype"`
	TimeStamp *TimeStamp `json:"ts"`
	UserID    UserID     `json:"user"`
	Text      string     `json:"text"`
	ItemType  ItemType   `json:"item_type"`
	Item      *Item      `json:"item"`
}

type PinAdded struct {
	TypedEvent
	UserID    UserID     `json:"user"`
	ChannelID ChannelID  `json:"channel_id"`
	TimeStamp *TimeStamp `json:"event_ts"`
	Item      *Item      `json:"item"`
}

type PinRemoved struct {
	TypedEvent
	UserID    UserID     `json:"user"`
	ChannelID ChannelID  `json:"channel_id"`
	Item      *Item      `json:"item"`
	HasPins   bool       `json:"has_pins"`
	TimeStamp *TimeStamp `json:"event_ts"`
}

type PreferenceChanged struct {
	TypedEvent
	Name  string `json:"name"`
	Value string `json:"value"`
}

type PresenceChanged struct {
	TypedEvent
	UserID   UserID `json:"user"`
	Presence string `json:"presence"`
}

type ReactionAdded struct {
	TypedEvent
	UserID      UserID     `json:"user"`
	Reaction    string     `json:"reaction"` // TODO actual value
	ItemOwnerID UserID     `json:"item_user"`
	Item        *Item      `json:"item"` // TODO message, file, file comment. only ids are given, right?
	TimeStamp   *TimeStamp `json:"event_ts"`
}

type ReactionRemoved struct {
	TypedEvent
	UserID      UserID     `json:"user"`
	Reaction    string     `json:"reaction"` // TODO actual value
	ItemOwnerID string     `json:"item_user"`
	Item        *Item      `json:"item"` // TODO message, file, file comment. only ids are given, right?
	TimeStamp   *TimeStamp `json:"event_ts"`
}

// ReconnectURL is currently unsupported and experimental
// https://api.slack.com/events/reconnect_url
type ReconnectURL struct {
	TypedEvent
}

// ResourcesAdded is delivered as users install your Slack app, add your app to channels and conversations,
// or approve your app for additional permissions and resources.
// https://api.slack.com/events/resources_added
type ResourcesAdded struct {
	TypedEvent
	Resources []*struct {
		Resource *Resource `json:"resource"`
		Scopes   []string  `json:"scopes"`
	} `json:"resources"`
}

// ResourcesRemoved is delivered as users uninstall your Slack app and remove your app to channels & conversations.
// https://api.slack.com/events/resources_removed
type ResourcesRemoved struct {
	TypedEvent
	Resources []*struct {
		Resource *Resource `json:"resource"`
		Scopes   []string  `json:"scopes"`
	} `json:"resources"`
}

// ScopeDenied is sent when scope you requested is denied.
// https://api.slack.com/events/scope_denied
type ScopeDenied struct {
	TypedEvent
	Scopes    []string `json:"scopes"`
	TriggerID string   `json:"trigger_id"`
}

// ScopeGranted is sent when scope you requested is granted.
// https://api.slack.com/events/scope_granted
type ScopeGranted struct {
	TypedEvent
	Scopes    []string `json:"scopes"`
	TriggerID string   `json:"trigger_id"`
}

type StarAdded struct {
	TypedEvent
	UserID    UserID     `json:"user"`
	Item      *Item      `json:"item"` // TODO message, file, file comment. only ids are given, right?
	TimeStamp *TimeStamp `json:"event_ts"`
}

type StarRemoved struct {
	TypedEvent
	UserID    UserID     `json:"user"`
	Item      *Item      `json:"item"` // TODO message, file, file comment. only ids are given, right?
	TimeStamp *TimeStamp `json:"event_ts"`
}

type SubTeamCreated struct {
	TypedEvent
	SubTeam *SubTeam `json:"subteam"`
}

type SubTeamMembersChanged struct {
	TypedEvent
	SubTeamID         SubTeamID  `json:"subteam_id"`
	TeamID            TeamID     `json:"team_id"`
	PreviousUpdate    *TimeStamp `json:"date_previous_update"`
	Updated           *TimeStamp `json:"date_update"`
	AddedUserIDs      []UserID   `json:"added_users"`
	AddedUserCount    int        `json:"added_users_count,string"`
	RemovedUserIDs    []UserID   `json:"removed_users"`
	RemovedUsersCount int        `json:"removed_users_count,string"`
}

type SubTeamSelfAdded struct {
	TypedEvent
	SubTeamID SubTeamID `json:"subteam_id"`
}

type SubTeamSelfRemoved struct {
	TypedEvent
	SubTeamID SubTeamID `json:"subteam_id"`
}

type SubTeamUpdated struct {
	TypedEvent
	SubTeam *SubTeam `json:"subteam"`
}

type TeamDomainChanged struct {
	TypedEvent
	URL    string `json:"url"`
	Domain string `json:"domain"`
}

type TeamJoined struct {
	TypedEvent
	User *User `json:"user"`
}

// TeamMigrationStarted is sent when chat group is migrated between servers.
// "The WebSocket connection will close immediately after it is sent.
// *snip* By the time a client has reconnected the process is usually complete, so the impact is minimal."
// https://api.slack.com/events/team_migration_started
type TeamMigrationStarted struct {
	TypedEvent
}

type TeamPlanChanged struct {
	TypedEvent
	Plan         string   `json:"plan"` // currently "", std, and plus
	CanAddUra    bool     `json:"can_add_ura"`
	PaidFeatures []string `json:"paid_features"`
}

type TeamPreferenceChanged struct {
	TypedEvent
	Name  string `json:"name"`
	Value bool   `json:"value"`
}

type TeamProfileChanged struct {
	TypedEvent
	Profile *struct {
		// https://api.slack.com/events/team_profile_change
		// TODO: Only the modified field definitions are included in the payload.
		Fields []*struct {
			ID string `json:"id"`
		} `json:"fields"`
	} `json:"profile"`
}

type TeamProfileDeleted struct {
	TypedEvent
	Profile *struct {
		Fields []string `json:"fields"`
	} `json:"profile"`
}

type TeamProfileReordered struct {
	TypedEvent
	Profile *struct {
		Fields []*struct {
			ID    string `json:"id"`
			Order int    `json:"ordering"`
		} `json:"fields"`
	} `json:"profile"`
}

type TeamRenamed struct {
	TypedEvent
	Name string `json:"name"`
}

// TokensRevoked is sent when API tokens are revoked.
// https://api.slack.com/events/tokens_revoked
type TokensRevoked struct {
	TypedEvent
	Tokens *struct {
		OAuth []string `json:"oauth"`
		Bot   []string `json:"bot"`
	}
}

type UserChanged struct {
	TypedEvent
	User *User `json:"user"`
}

// UserResourceDenied is sent when a user declines to grant your workspace app the permissions you recently requested with "apps.permissions.users.request."
// https://api.slack.com/events/user_resource_denied
type UserResourceDenied struct {
	TypedEvent
	UserID    UserID   `json:"user"`
	Scopes    []string `json:"scopes"`
	TriggerID string   `json:"trigger_id"`
}

// UserResourceGranted is sent when a user grant your workspace app the permissions you recently requested with "apps.permissions.users.request."
// https://api.slack.com/events/user_resource_granted
type UserResourceGranted struct {
	TypedEvent
	UserID    UserID   `json:"user"`
	Scopes    []string `json:"scopes"`
	TriggerID string   `json:"trigger_id"`
}

// UserResourceRemoved is sent when a user removes an existing grant for your workspace app.
// https://api.slack.com/events/user_resource_removed
type UserResourceRemoved struct {
	TypedEvent
	UserID    UserID `json:"user"`
	TriggerID string `json:"trigger_id"`
}

type UserTyping struct {
	TypedEvent
	ChannelID ChannelID `json:"channel"`
	UserID    UserID    `json:"user"`
}

// Below comes some types that are commonly shared amongst multiple events

type App struct {
	ID                     AppID  `json:"id"`
	Name                   string `json:"name"`
	Description            string `json:"description"`
	HelpURL                string `json:"help_url"`
	PrivacyPolicyURL       string `json:"privacy_policy_url"`
	AppHomepageURL         string `json:"app_homepage_url"`
	AppDirectoryURL        string `json:"app_directory_url"`
	IsAppDirectoryApproved bool   `json:"is_app_directory_approved"`
	IsInternal             bool   `json:"is_internal"`
	AdditionalInfo         string `json:"additional_info"`
}

type AppRequest struct {
	ID                 string `json:"id"`
	App                *App   `json:"app"`
	PreviousResolution *struct {
		Status string             `json:"status"`
		Scopes []*AppRequestScope `json:"scopes"`
	} `json:"previous_resolution"`
	User *struct {
		ID    UserID `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"user"`
	Team *struct {
		ID     TeamID `json:"id"`
		Name   string `json:"name"`
		Domain string `json:"domain"`
	} `json:"team"`
	Scopes  []*AppRequestScope `json:"scopes"`
	Message string             `json:"message"`
}

type AppRequestScope struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsSensitive bool   `json:"is_sensitive"`
	TokenType   string `json:"token_type"`
}

type Bot struct {
	ID    BotID    `json:"id"`
	AppID AppID    `json:"app_id"`
	Name  string   `json:"name"`
	Icon  *BotIcon `json:"icons"`
}

type BotIcon struct {
	Image48 string `json:"image_48"`
}

type ChangedHistory struct {
	Latest         *TimeStamp `json:"latest"`
	TimeStamp      *TimeStamp `json:"ts"`
	EventTimeStamp *TimeStamp `json:"event_ts"`
}

type Comment struct {
	ID        CommentID  `json:"id"`
	TimeStamp *TimeStamp `json:"created"`
	UserID    UserID     `json:"user"`
	Content   string     `json:"comment"`
}

// https://api.slack.com/types/file
// Most "file" object in rtm events only contains id field.
type File struct {
	ID                 FileID      `json:"id"`
	TimeStamp          *TimeStamp  `json:"created"`
	Name               string      `json:"name"`
	Title              string      `json:"title"`
	MimeType           string      `json:"mimetype"`
	FileType           string      `json:"filetype"`
	PrettyType         string      `json:"pretty_type"`
	UserID             UserID      `json:"user"`
	Mode               string      `json:"mode"`
	Editable           bool        `json:"editable"`
	IsExternal         bool        `json:"is_external"`
	ExternalType       string      `json:"external_type"`
	UserName           string      `json:"username"`
	Size               int         `json:"size"`
	URLPrivate         string      `json:"url_private"`
	URLPrivateDownload string      `json:"url_private_download"`
	Thumb64            string      `json:"thumb_64"`
	Thumb80            string      `json:"thumb_80"`
	Thumb360           string      `json:"thumb_360"`
	Thumb360Gif        string      `json:"thumb_360_gif"`
	Thumb360W          int         `json:"thumb_360_w"`
	Thumb360H          int         `json:"thumb_360_h"`
	Thumb480           string      `json:"thumb_480"`
	Thumb480W          int         `json:"thumb_480_w"`
	Thumb480H          int         `json:"thumb_480_h"`
	Thumb160           string      `json:"thumb_160"`
	Permalink          string      `json:"permalink"`
	PermalinkPublic    string      `json:"permalink_public"`
	EditLink           string      `json:"edit_link"`
	Preview            string      `json:"preview"`
	PreviewHighlight   string      `json:"preview_highlight"`
	Lines              int         `json:"lines"`
	LinesMore          int         `json:"lines_more"`
	IsPublic           bool        `json:"is_public"`
	PublicURLShared    bool        `json:"public_url_shared"`
	DisplayAsBot       bool        `json:"display_as_bot"`
	ChannelIDs         []ChannelID `json:"channels"`
	Groups             []string    `json:"groups"`
	Ims                []string    `json:"ims"`
	InitialComment     *Comment    `json:"initial_comment"`
	NumStars           int         `json:"num_stars"`
	IsStarred          bool        `json:"is_starred"`
	PinnedTo           []string    `json:"pinned_to"`
	Reactions          []*struct {
		Name    string   `json:"name"`
		Count   int      `json:"count"`
		UserIDs []UserID `json:"users"`
	} `json:"reactions"`
	CommentsCount int `json:"comments_count"`
}

// Item can be any object with type of Message, File, or Comment.
type Item struct {
	Type      string     `json:"type"`
	ChannelID ChannelID  `json:"channel"`
	Message   *Message   `json:"message"`
	File      *File      `json:"file"`
	Comment   *Comment   `json:"comment"`
	TimeStamp *TimeStamp `json:"ts"`
}

type MarkedAsRead struct {
	ChannelID ChannelID  `json:"channel"`
	TimeStamp *TimeStamp `json:"ts"`
}

type Resource struct {
	Type  string `json:"type"`
	Grant *struct {
		Type       string     `json:"type"`
		ResourceID ResourceID `json:"resource_id"`
	} `json:"grant"`
}

type SubTeam struct {
	ID          SubTeamID  `json:"id"`
	TeamID      TeamID     `json:"team_id"`
	IsUserGroup bool       `json:"is_usergroup"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Handle      string     `json:"handle"`
	IsExternal  bool       `json:"is_external"`
	Created     *TimeStamp `json:"date_create"`
	Updated     *TimeStamp `json:"date_update"`
	Deleted     *TimeStamp `json:"date_delete"`
	AutoType    string     `json:"auto_type"`
	CreatorID   UserID     `json:"created_by"`
	UpdatedBy   string     `json:"updated_by"`
	UserCount   int        `json:"user_count,string"`
	UserIDs     []UserID   `json:"users"`
}

type User struct {
	ID       UserID `json:"id"`
	Name     string `json:"name"`
	Deleted  bool   `json:"deleted"`
	Color    string `json:"color"`
	RealName string `json:"real_name"`
	TZ       string `json:"tz"`
	TZLabel  string `json:"tz_label"`
	TZOffset int    `json:"tz_offset"`
	Profile  *struct {
		FirstName          string `json:"first_name"`
		LastName           string `json:"last_name"`
		RealName           string `json:"real_name"`
		RealNameNormalized string `json:"real_name_normalized"`
		Email              string `json:"email"`
		Skype              string `json:"skype"`
		Phone              string `json:"phone"`
		Image24            string `json:"image_24"`
		Image32            string `json:"image_32"`
		Image48            string `json:"image_48"`
		Image72            string `json:"image_72"`
		Image192           string `json:"image_192"`
		ImageOriginal      string `json:"image_original"`
		Title              string `json:"title"`
	} `json:"profile"`
	IsBot             bool   `json:"is_bot"`
	IsAdmin           bool   `json:"is_admin"`
	IsOwner           bool   `json:"is_owner"`
	IsPrimaryOwner    bool   `json:"is_primary_owner"`
	IsRestricted      bool   `json:"is_restricted"`
	IsUltraRestricted bool   `json:"is_ultra_restricted"`
	Has2FA            bool   `json:"has_2fa"`
	HasFiles          bool   `json:"has_files"`
	Presence          string `json:"presence"`
}

type View struct {
	ID                 ViewID     `json:"id"`
	TeamID             TeamID     `json:"team_id"`
	Type               string     `json:"type"`
	Blocks             []Block    `json:"blocks"`
	PrivateMetadata    string     `json:"private_metadata"`
	CallbackID         string     `json:"callback_id"`
	State              *ViewState `json:"state"`
	Hash               string     `json:"hash"`
	ClearOnClose       bool       `json:"clear_on_close"`
	NotifyOnClose      bool       `json:"notify_on_close"`
	RootViewID         string     `json:"root_view_id"`
	AppID              AppID      `json:"app_id"`
	ExternalID         string     `json:"external_id"`
	AppInstalledTeamID TeamID     `json:"app_installed_team_id"`
	BotID              BotID      `json:"bot_id"`
}

func (v *View) UnmarshalJSON(b []byte) error {
	type alias View
	t := &struct {
		*alias
		Blocks []json.RawMessage `json:"blocks"`
	}{
		alias: (*alias)(v),
	}
	err := json.Unmarshal(b, t)
	if err != nil {
		return err
	}

	var blocks []Block
	for _, elem := range t.Blocks {
		block, err := UnmarshalBlock(elem)
		if err != nil {
			return xerrors.Errorf("failed to unmarshal given block: %w", err)
		}
		blocks = append(blocks, block)
	}
	v.Blocks = blocks
	return nil
}

type ViewState struct {
	Values map[BlockID]map[ActionID]*ViewStateValue `json:"values"` // https://api.slack.com/reference/interaction-payloads/views#view_submission_fields
}

type ViewStateValue struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
