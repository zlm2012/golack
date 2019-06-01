package rtmapi

import (
	"strings"
)

// EventType represents the type of event sent from slack.
// Event is passed to client in a form of JSON string, which has a field named "type."
type EventType string

// EventTyper represents an RTM API payload with "type."
type EventTyper interface {
	EventType() EventType
}

// List of available EventTypes
const (
	UnsupportedEvent           EventType = "unsupported"
	AccountsChangedEvent                 = "accounts_changed"
	BotAddedEvent                        = "bot_added"
	BotChangedEvent                      = "bot_changed"
	ChannelArchivedEvent                 = "channel_archive"
	ChannelCreatedEvent                  = "channel_created"
	ChannelDeletedEvent                  = "channel_deleted"
	ChannelHistoryChangedEvent           = "channel_history_changed"
	ChannelJoinedEvent                   = "channel_joined"
	ChannelLeftEvent                     = "channel_left"
	ChannelMarkedEvent                   = "channel_marked"
	ChannelRenameEvent                   = "channel_rename"
	ChannelUnarchiveEvent                = "channel_unarchive"
	CommandsChangedEvent                 = "commands_changed"
	DNDUpdatedEvent                      = "dnd_updated"
	DNDUpdatedUserEvent                  = "dnd_updated_user"
	EmailDomainChangedEvent              = "email_domain_changed"
	EmojiChangedEvent                    = "emoji_changed"
	FileChangeEvent                      = "file_change"
	FileCommentAddedEvent                = "file_comment_added"
	FileCommentDeletedEvent              = "file_comment_deleted"
	FileCommentEditedEvent               = "file_comment_edited"
	FileCreatedEvent                     = "file_created"
	FileDeletedEvent                     = "file_deleted"
	FilePublicEvent                      = "file_public"
	FileSharedEvent                      = "file_shared"
	FileUnsharedEvent                    = "file_unshared"
	GoodByeEvent                         = "goodbye"
	GroupArchiveEvent                    = "group_archive"
	GroupCloseEvent                      = "group_close"
	GroupDeletedEvent                    = "group_deleted"
	GroupHistoryChangedEvent             = "group_history_changed"
	GroupJoinedEvent                     = "group_joined"
	GroupLeftEvent                       = "group_left"
	GroupMarkedEvent                     = "group_marked"
	GroupOpenEvent                       = "group_open"
	GroupRenameEvent                     = "group_rename"
	GroupUnarchiveEvent                  = "group_unarchive"
	HelloEvent                           = "hello"
	IMCloseEvent                         = "im_close"
	IMCreatedEvent                       = "im_created"
	IMHistoryChangedEvent                = "im_history_changed"
	IMMarkedEvent                        = "im_marked"
	IMOpenEvent                          = "im_open"
	ManualPresenceChangeEvent            = "manual_presence_change"
	MemberJoinedChannelEvent             = "member_joined_channel"
	MessageEvent                         = "message"
	PinAddedEvent                        = "pin_added"
	PinRemovedEvent                      = "pin_removed"
	PrefChangeEvent                      = "pref_change"
	PresenceChangeEvent                  = "presence_change"
	PresenceQueryEvent                   = "presence_query"
	PresenceSubEvent                     = "presence_sub"
	ReactionAddedEvent                   = "reaction_added"
	ReactionRemovedEvent                 = "reaction_removed"
	ReconnectURLEvent                    = "reconnect_url"
	StarAddedEvent                       = "star_added"
	StarRemovedEvent                     = "star_removed"
	SubTeamCreatedEvent                  = "subteam_created"
	SubTeamMembersChangedEvent           = "subteam_members_changed"
	SubTeamSelfAddedEvent                = "subteam_self_added"
	SubTeamSelfRemovedEvent              = "subteam_self_removed"
	SubTeamUpdatedEvent                  = "subteam_updated"
	TeamDomainChangeEvent                = "team_domain_change"
	TeamJoinEvent                        = "team_join"
	TeamMigrationStartedEvent            = "team_migration_started"
	TeamPlanChangedEvent                 = "team_plan_change"
	TeamPrefChangeEvent                  = "team_pref_change"
	TeamProfileChangeEvent               = "team_profile_change"
	TeamProfileDeleteEvent               = "team_profile_delete"
	TeamProfileReorderEvent              = "team_profile_reorder"
	TeamRenameEvent                      = "team_rename"
	UserChangeEvent                      = "user_change"
	UserTypingEvent                      = "user_typing"
	PingEvent                            = "ping"
	PongEvent                            = "pong"
)

var (
	possibleEvents = [...]EventType{
		UnsupportedEvent, AccountsChangedEvent, BotAddedEvent, BotChangedEvent, ChannelArchivedEvent, ChannelCreatedEvent, ChannelDeletedEvent,
		ChannelHistoryChangedEvent, ChannelJoinedEvent, ChannelLeftEvent, ChannelMarkedEvent, ChannelRenameEvent, ChannelUnarchiveEvent,
		CommandsChangedEvent, DNDUpdatedEvent, DNDUpdatedUserEvent, EmailDomainChangedEvent, EmojiChangedEvent, FileChangeEvent,
		FileCommentAddedEvent, FileCommentDeletedEvent, FileCommentEditedEvent, FileCreatedEvent, FileDeletedEvent, FilePublicEvent,
		FileSharedEvent, FileUnsharedEvent, GoodByeEvent, GroupArchiveEvent, GroupCloseEvent, GroupDeletedEvent, GroupHistoryChangedEvent, GroupJoinedEvent,
		GroupLeftEvent, GroupMarkedEvent, GroupOpenEvent, GroupRenameEvent, GroupUnarchiveEvent, HelloEvent, IMCloseEvent, IMCreatedEvent,
		IMHistoryChangedEvent, IMMarkedEvent, IMOpenEvent, ManualPresenceChangeEvent, MemberJoinedChannelEvent, MessageEvent, PinAddedEvent, PinRemovedEvent,
		PrefChangeEvent, PresenceChangeEvent, PresenceQueryEvent, PresenceSubEvent, ReactionAddedEvent, ReactionRemovedEvent, ReconnectURLEvent, StarAddedEvent, StarRemovedEvent,
		SubTeamCreatedEvent, SubTeamMembersChangedEvent, SubTeamSelfAddedEvent, SubTeamSelfRemovedEvent, SubTeamUpdatedEvent, TeamDomainChangeEvent, TeamJoinEvent,
		TeamMigrationStartedEvent, TeamPlanChangedEvent, TeamPrefChangeEvent, TeamProfileChangeEvent, TeamProfileDeleteEvent, TeamProfileReorderEvent,
		TeamRenameEvent, UserChangeEvent, UserTypingEvent, PingEvent, PongEvent,
	}
)

// UnmarshalText parses a given event value to EventType.
// This method is mainly used by encode/json.
func (eventType *EventType) UnmarshalText(b []byte) error {
	*eventType = AtoEventType(string(b))
	return nil
}

// AtoEventType converts given string to corresponding EventType.
func AtoEventType(str string) EventType {
	for _, val := range possibleEvents {
		if str == val.String() {
			return val
		}
	}

	return UnsupportedEvent
}

// String returns the stringified event name, which corresponds to the one sent from/to slack RTM endpoint.
func (eventType EventType) String() string {
	return string(eventType)
}

// MarshalText returns the stringified value of slack event.
// This method is mainly used by encode/json.
func (eventType *EventType) MarshalText() ([]byte, error) {
	str := eventType.String()

	if strings.Compare(str, "") == 0 {
		return []byte(UnsupportedEvent), nil
	}

	return []byte(str), nil
}

// TypedEvent takes care of events that have "type" field in its JSON representation.
// The API document, https://api.slack.com/rtm#events, states as follows:
// "Every event has a type property which describes the type of event."
type TypedEvent struct {
	Type EventType `json:"type,omitempty"`
}

var _ EventTyper = TypedEvent{}

func (te TypedEvent) EventType() EventType {
	return te.Type
}

// SubType may given as a part of message payload to describe detailed content.
type SubType string

const (
	// List of available SubTypes
	Empty            SubType = "" // can be absent
	BotMessage               = "bot_message"
	ChannelArchive           = "channel_archive"
	ChannelJoin              = "channel_join"
	ChannelLeave             = "channel_leave"
	ChannelName              = "channel_name"
	ChannelPurpose           = "channel_purpose"
	ChannelTopic             = "channel_topic"
	ChannelUnarchive         = "channel_unarchive"
	FileComment              = "file_comment"
	FileMention              = "file_mention"
	FileShare                = "file_share"
	GroupArchive             = "group_archive"
	GroupJoin                = "group_join"
	GroupLeave               = "group_leave"
	GroupName                = "group_name"
	GroupPurpose             = "group_purpose"
	GroupTopic               = "group_topic"
	GroupUnarchive           = "group_unarchive"
	MeMessage                = "me_message"
	MessageChanged           = "message_changed"
	MessageDeleted           = "message_deleted"
	PinnedItem               = "pinned_item"
	UnpinnedItem             = "unpinned_item"
)

var (
	possibleSubTypes = [...]SubType{
		BotMessage, ChannelArchive, ChannelJoin, ChannelLeave, ChannelName, ChannelPurpose, ChannelTopic,
		ChannelUnarchive, FileComment, FileMention, FileShare, GroupArchive, GroupJoin, GroupLeave,
		GroupName, GroupPurpose, GroupTopic, GroupUnarchive, MeMessage, MessageChanged, MessageDeleted,
		PinnedItem, UnpinnedItem,
	}
)

// UnmarshalText parses a given subtype value to SubType
// This method is mainly used by encode/json.
func (subType *SubType) UnmarshalText(b []byte) error {
	*subType = AtoSubType(string(b))
	return nil
}

// AtoSubType converts given string to corresponding SubType.
func AtoSubType(str string) SubType {
	for _, val := range possibleSubTypes {
		if str == val.String() {
			return val
		}
	}

	return Empty
}

// String returns the stringified subtype name, which corresponds to the one sent from/to slack RTM endpoint.
func (subType SubType) String() string {
	return string(subType)
}

// MarshalText returns the stringified value of slack subtype.
// This method is mainly used by encode/json.
func (subType *SubType) MarshalText() ([]byte, error) {
	str := subType.String()

	if strings.Compare(str, "") == 0 {
		return []byte(""), nil // EMPTY
	}

	return []byte(str), nil
}

// CommonMessage contains some common fields of message event.
// See SubType field to distinguish corresponding event struct.
// https://api.slack.com/events/message#message_subtypes
type CommonMessage struct {
	TypedEvent

	// Regular user message and some miscellaneous message share the common type of "message."
	// So take a look at subtype to distinguish. Regular user message has empty subtype.
	SubType SubType `json:"subtype"`
}
