package rtmapi

import "github.com/oklahomer/golack/slackobject"

// AccountsChanged event indicates that the list of accounts a user is signed into has changed.
// https://api.slack.com/events/accounts_changed
type AccountsChanged struct {
	TypedEvent
}

type Bot struct {
	ID    slackobject.BotID `json:"id"`
	AppID slackobject.AppID `json:"app_id"`
	Name  string            `json:"name"`
	Icon  *struct {
		Image48 string `json:"image_48"`
	} `json:"icons"`
}

type BotAdded struct {
	TypedEvent
	Bot *Bot `json:"bot"`
}

type BotChanged struct {
	TypedEvent
	Bot *Bot `json:"bot"`
}

type ChannelArchived struct {
	TypedEvent
	ChannelID slackobject.ChannelID `json:"channel"`
	UserID    slackobject.UserID    `json:"user"`
}

type ChannelCreated struct {
	TypedEvent
	Channel *struct {
		ID        slackobject.ChannelID `json:"id"`
		Name      string                `json:"name"`
		Created   *TimeStamp            `json:"created"`
		CreatorID slackobject.UserID    `json:"creator"`
	} `json:"channel"`
}

type ChannelDeleted struct {
	TypedEvent
	ChannelID slackobject.ChannelID `json:"channel"`
}

type historyChangedEvent struct {
	Latest         *TimeStamp `json:"latest"`
	TimeStamp      *TimeStamp `json:"ts"`
	EventTimeStamp *TimeStamp `json:"event_ts"`
}

type ChannelHistoryChanged struct {
	TypedEvent
	historyChangedEvent
}

type ChannelJoined struct {
	TypedEvent
	ChannelID slackobject.ChannelID `json:"channel"`
}

type ChannelLeft struct {
	TypedEvent
	ChannelID slackobject.ChannelID `json:"channel"`
}

type markedAsReadEvent struct {
	ChannelID slackobject.ChannelID `json:"channel"`
	TimeStamp *TimeStamp            `json:"ts"`
}

type ChannelMarked struct {
	TypedEvent
	markedAsReadEvent
}

type ChannelRenamed struct {
	TypedEvent
	Channel *struct {
		ID      slackobject.ChannelID `json:"id"`
		Name    string                `json:"name"`
		Created *TimeStamp            `json:"created"`
	} `json:"channel"`
}

type ChannelUnarchived struct {
	TypedEvent
	ChannelID slackobject.ChannelID `json:"channel"`
	UserID    slackobject.UserID    `json:"user"`
}

type CommandsChanged struct {
	TypedEvent
	TimeStamp *TimeStamp `json:"event_ts"`
}

type DNDUpdated struct {
	TypedEvent
	UserID    slackobject.UserID `json:"user"`
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
	UserID    slackobject.UserID `json:"user"`
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

type Comment struct {
	ID      slackobject.CommentID `json:"id"`
	Created *TimeStamp            `json:"created"`
	UserID  slackobject.UserID    `json:"user"`
	Content string                `json:"comment"`
}

// https://api.slack.com/types/file
// TODO See if this object should stay in rtmapi package.
// It seems like to be used in REST API.
// Most "file" object in rtm events only contains id field.
type File struct {
	ID                 slackobject.FileID      `json:"id"`
	TimeStamp          *TimeStamp              `json:"created"`
	Name               string                  `json:"name"`
	Title              string                  `json:"title"`
	MimeType           string                  `json:"mimetype"`
	FileType           string                  `json:"filetype"`
	PrettyType         string                  `json:"pretty_type"`
	UserID             slackobject.UserID      `json:"user"`
	Mode               string                  `json:"mode"`
	Editable           bool                    `json:"editable"`
	IsExternal         bool                    `json:"is_external"`
	ExternalType       string                  `json:"external_type"`
	UserName           string                  `json:"username"`
	Size               int                     `json:"size"`
	URLPrivate         string                  `json:"url_private"`
	URLPrivateDownload string                  `json:"url_private_download"`
	Thumb64            string                  `json:"thumb_64"`
	Thumb80            string                  `json:"thumb_80"`
	Thumb360           string                  `json:"thumb_360"`
	Thumb360Gif        string                  `json:"thumb_360_gif"`
	Thumb360W          int                     `json:"thumb_360_w"`
	Thumb360H          int                     `json:"thumb_360_h"`
	Thumb480           string                  `json:"thumb_480"`
	Thumb480W          int                     `json:"thumb_480_w"`
	Thumb480H          int                     `json:"thumb_480_h"`
	Thumb160           string                  `json:"thumb_160"`
	Permalink          string                  `json:"permalink"`
	PermalinkPublic    string                  `json:"permalink_public"`
	EditLink           string                  `json:"edit_link"`
	Preview            string                  `json:"preview"`
	PreviewHighlight   string                  `json:"preview_highlight"`
	Lines              int                     `json:"lines"`
	LinesMore          int                     `json:"lines_more"`
	IsPublic           bool                    `json:"is_public"`
	PublicURLShared    bool                    `json:"public_url_shared"`
	DisplayAsBot       bool                    `json:"display_as_bot"`
	ChannelIDs         []slackobject.ChannelID `json:"channels"`
	Groups             []string                `json:"groups"`
	Ims                []string                `json:"ims"`
	InitialComment     *Comment                `json:"initial_comment"`
	NumStars           int                     `json:"num_stars"`
	IsStarred          bool                    `json:"is_starred"`
	PinnedTo           []string                `json:"pinned_to"`
	Reactions          []*struct {
		Name    string               `json:"name"`
		Count   int                  `json:"count"`
		UserIDs []slackobject.UserID `json:"users"`
	} `json:"reactions"`
	CommentsCount int `json:"comments_count"`
}

type FileChanged struct {
	TypedEvent
	FileID slackobject.FileID `json:"file_id"`
	File   *struct {
		// https://api.slack.com/events/file_change
		// This value is identical to FileID field
		ID slackobject.FileID `json:"id"`
	} `json:"file"`
}

type FileCommentAdded struct {
	TypedEvent
	FileID slackobject.FileID `json:"file_id"`
	File   *struct {
		// https://api.slack.com/events/file_change
		// This value is identical to FileID field
		ID slackobject.FileID `json:"id"`
	} `json:"file"`
	Comment *Comment `json:"comment"`
}

type FileCommentDeleted struct {
	TypedEvent
	FileID slackobject.FileID `json:"file_id"`
	File   *struct {
		// https://api.slack.com/events/file_change
		// This value is identical to FileID field
		ID slackobject.FileID `json:"id"`
	} `json:"file"`
	CommentID string `json:"comment"`
}

type FileCommentEdited struct {
	TypedEvent
	FileID slackobject.FileID `json:"file_id"`
	File   *struct {
		// https://api.slack.com/events/file_change
		// This value is identical to FileID field
		ID slackobject.FileID `json:"id"`
	} `json:"file"`
	Comment *Comment `json:"comment"`
}

type FileCreated struct {
	TypedEvent
	FileID slackobject.FileID `json:"file_id"`
	File   *struct {
		// https://api.slack.com/events/file_change
		// This value is identical to FileID field
		ID slackobject.FileID `json:"id"`
	} `json:"file"`
}

type FileDeleted struct {
	TypedEvent
	FileID    slackobject.FileID `json:"file_id"`
	TimeStamp *TimeStamp         `json:"event_ts"`
}

type FilePublished struct {
	TypedEvent
	FileID slackobject.FileID `json:"file_id"`
	File   *struct {
		// https://api.slack.com/events/file_change
		// This value is identical to FileID field
		ID slackobject.FileID `json:"id"`
	} `json:"file"`
}

type FileShared struct {
	TypedEvent
	FileID slackobject.FileID `json:"file_id"`
	File   *struct {
		// https://api.slack.com/events/file_change
		// This value is identical to FileID field
		ID slackobject.FileID `json:"id"`
	} `json:"file"`
}

type FileUnshared struct {
	TypedEvent
	FileID slackobject.FileID `json:"file_id"`
	File   *struct {
		// https://api.slack.com/events/file_change
		// This value is identical to FileID field
		ID slackobject.FileID `json:"id"`
	} `json:"file"`
}

type GoodBye struct {
	TypedEvent
}

type GroupArchived struct {
	TypedEvent
	ChannelID slackobject.ChannelID `json:"channel"`
}

type GroupClosed struct {
	TypedEvent
	UserID    slackobject.UserID    `json:"user"`
	ChannelID slackobject.ChannelID `json:"channel"`
}

type GroupDeleted struct {
	TypedEvent
	ChannelID slackobject.ChannelID `json:"channel"`
}

type GroupHistoryChanged struct {
	TypedEvent
	historyChangedEvent
}

type GroupJoined struct {
	TypedEvent
	ChannelID slackobject.ChannelID `json:"channel"`
}

type GroupLeft struct {
	TypedEvent
	ChannelID slackobject.ChannelID `json:"channel"`
}

type GroupMarked struct {
	TypedEvent
	markedAsReadEvent
}

type GroupOpened struct {
	TypedEvent
	UserID    slackobject.UserID    `json:"user"`
	ChannelID slackobject.ChannelID `json:"channel"`
}

type GroupRenamed struct {
	TypedEvent
	Channel *struct {
		ID      slackobject.ChannelID `json:"id"`
		Name    string                `json:"name"`
		Created *TimeStamp            `json:"created"`
	} `json:"channel"`
}

type GroupUnarchived struct {
	TypedEvent
	ChannelID slackobject.ChannelID `json:"channel"`
}

// Hello event is sent from slack when WebSocket connection is successfully established.
// https://api.slack.com/events/hello
type Hello struct {
	TypedEvent
}

type IMClosed struct {
	TypedEvent
	UserID    slackobject.UserID    `json:"user"`
	ChannelID slackobject.ChannelID `json:"channel"`
}

type IMCreated struct {
	TypedEvent
	UserID  slackobject.UserID `json:"user"`
	Channel *struct {
		ID slackobject.ChannelID `json:"id"`
	} `json:"channel"`
}

type IMHistoryChanged struct {
	TypedEvent
	historyChangedEvent
}

type IMMarked struct {
	TypedEvent
	markedAsReadEvent
}

type IMOpened struct {
	TypedEvent
	UserID    slackobject.UserID    `json:"user"`
	ChannelID slackobject.ChannelID `json:"channel"`
}

type PresenceManuallyChanged struct {
	TypedEvent
	Presence string `json:"presence"` // TODO Actual values other than "away"
}

type MemberJoinedChannel struct {
	TypedEvent
	UserID      slackobject.UserID    `json:"user"`
	ChannelID   slackobject.ChannelID `json:"channel"`
	ChannelType string                `json:"channel_type"` // C or G. ref. https://api.slack.com/events/member_joined_channel
	TeamID      slackobject.TeamID    `json:"team"`
	InviterID   slackobject.UserID    `json:"inviter"` // Empty when the user joins by herself. ref. https://api.slack.com/events/member_joined_channel
}

// Message represent message event on RTM.
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
	ChannelID       slackobject.ChannelID `json:"channel"`
	SenderID        slackobject.UserID    `json:"user"`
	Text            string                `json:"text"`
	TimeStamp       *TimeStamp            `json:"ts"`
	ThreadTimeStamp *TimeStamp            `json:"thread_ts"` // https://api.slack.com/docs/message-threading
}

// Item can be any object with type of Message, File, or Comment.
type Item struct {
	Type      string                `json:"type"`
	ChannelID slackobject.ChannelID `json:"channel"`
	Message   *Message              `json:"message"`
	File      *File                 `json:"file"`
	Comment   *Comment              `json:"comment"`
	TimeStamp *TimeStamp            `json:"ts"`
}

type PinAdded struct {
	TypedEvent
	UserID    slackobject.UserID    `json:"user"`
	ChannelID slackobject.ChannelID `json:"channel_id"`
	TimeStamp *TimeStamp            `json:"event_ts"`
	Item      *Item                 `json:"item"`
}

type PinRemoved struct {
	TypedEvent
	UserID    slackobject.UserID    `json:"user"`
	ChannelID slackobject.ChannelID `json:"channel_id"`
	Item      *Item                 `json:"item"`
	HasPins   bool                  `json:"has_pins"`
	TimeStamp *TimeStamp            `json:"event_ts"`
}

type PreferenceChanged struct {
	TypedEvent
	Name  string `json:"name"`
	Value string `json:"value"`
}

type PresenceChange struct {
	TypedEvent
	UserID   slackobject.UserID `json:"user"`
	Presence string             `json:"presence"`
}

type PresenceQuery struct {
	TypedEvent
	UserIDs []slackobject.UserID `json:"ids"`
}

type PresenceSubscribe struct {
	TypedEvent
	UserIDs []slackobject.UserID `json:"ids"`
}

type ReactionAdded struct {
	TypedEvent
	UserID      slackobject.UserID `json:"user"`
	Reaction    string             `json:"reaction"` // TODO actual value
	ItemOwnerID slackobject.UserID `json:"item_user"`
	Item        *Item              `json:"item"` // TODO message, file, file comment. only ids are given, right?
	TimeStamp   *TimeStamp         `json:"event_ts"`
}

type ReactionRemoved struct {
	TypedEvent
	UserID      slackobject.UserID `json:"user"`
	Reaction    string             `json:"reaction"` // TODO actual value
	ItemOwnerID string             `json:"item_user"`
	Item        *Item              `json:"item"` // TODO message, file, file comment. only ids are given, right?
	TimeStamp   *TimeStamp         `json:"event_ts"`
}

// ReconnectURL is currently unsupported and experimental
// https://api.slack.com/events/reconnect_url
type ReconnectURL struct {
	TypedEvent
}

type StarAdded struct {
	TypedEvent
	UserID    slackobject.UserID `json:"user"`
	Item      *Item              `json:"item"` // TODO message, file, file comment. only ids are given, right?
	TimeStamp *TimeStamp         `json:"event_ts"`
}

type StarRemoved struct {
	TypedEvent
	UserID    slackobject.UserID `json:"user"`
	Item      *Item              `json:"item"` // TODO message, file, file comment. only ids are given, right?
	TimeStamp *TimeStamp         `json:"event_ts"`
}

type SubTeam struct {
	ID          slackobject.SubTeamID `json:"id"`
	TeamID      slackobject.TeamID    `json:"team_id"`
	IsUserGroup bool                  `json:"is_usergroup"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Handle      string                `json:"handle"`
	IsExternal  bool                  `json:"is_external"`
	Created     *TimeStamp            `json:"date_create"`
	Updated     *TimeStamp            `json:"date_update"`
	Deleted     *TimeStamp            `json:"date_delete"`
	AutoType    string                `json:"auto_type"`
	CreatorID   slackobject.UserID    `json:"created_by"`
	UpdatedBy   string                `json:"updated_by"`
	UserCount   int                   `json:"user_count,string"`
	UserIDs     []slackobject.UserID  `json:"users"`
}

type SubTeamCreated struct {
	TypedEvent
	SubTeam *SubTeam `json:"subteam"`
}

type SubTeamMembersChanged struct {
	TypedEvent
	SubTeamID         slackobject.SubTeamID `json:"subteam_id"`
	TeamID            slackobject.TeamID    `json:"team_id"`
	PreviousUpdate    *TimeStamp            `json:"date_previous_update"`
	Updated           *TimeStamp            `json:"date_update"`
	AddedUserIDs      []slackobject.UserID  `json:"added_users"`
	AddedUserCount    int                   `json:"added_users_count,string"`
	RemovedUserIDs    []slackobject.UserID  `json:"removed_users"`
	RemovedUsersCount int                   `json:"removed_users_count,string"`
}

type SubTeamSelfAdded struct {
	TypedEvent
	SubTeamID slackobject.SubTeamID `json:"subteam_id"`
}

type SubTeamSelfRemoved struct {
	TypedEvent
	SubTeamID slackobject.SubTeamID `json:"subteam_id"`
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

type User struct {
	ID       slackobject.UserID `json:"id"`
	Name     string             `json:"name"`
	Deleted  bool               `json:"deleted"`
	Color    string             `json:"color"`
	RealName string             `json:"real_name"`
	TZ       string             `json:"tz"`
	TZLabel  string             `json:"tz_label"`
	TZOffset int                `json:"tz_offset"`
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

type UserChanged struct {
	TypedEvent
	User *User `json:"user"`
}

type UserTyping struct {
	TypedEvent
	ChannelID slackobject.ChannelID `json:"channel"`
	UserID    slackobject.UserID    `json:"user"`
}

// Pong is given when client send Ping.
// https://api.slack.com/rtm#ping_and_pong
type Pong struct {
	TypedEvent
	ReplyTo uint `json:"reply_to"`
}

// MiscMessage represents some minor message events.
// TODO define each one with subtype field. This is just a representation of common subtyped payload
// https://api.slack.com/events/message#message_subtypes
type MiscMessage struct {
	CommonMessage
	TimeStamp *TimeStamp `json:"ts"`
}
