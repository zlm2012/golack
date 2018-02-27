package rtmapi

import "github.com/oklahomer/golack/slackobject"

// AccountsChanged event indicates that the list of accounts a user is signed into has changed.
// https://api.slack.com/events/accounts_changed
type AccountsChanged struct {
	CommonEvent
}

type Bot struct {
	ID    slackobject.BotID `json:"id"`
	AppID slackobject.AppID `json:"app_id"`
	Name  string            `json:"name"`
	Icons struct {
		Image48 string `json:"image_48"`
	} `json:"icons"`
}

type BotAdded struct {
	CommonEvent
	Bot *Bot `json:"bot"`
}

type BotChanged struct {
	CommonEvent
	Bot *Bot `json:"bot"`
}

type ChannelArchived struct {
	Type      string                `json:"type"`
	ChannelID slackobject.ChannelID `json:"channel"`
	User      slackobject.UserID    `json:"user"`
}

type ChannelCreated struct {
	CommonEvent
	Channel struct {
		ID      slackobject.ChannelID `json:"id"`
		Name    string                `json:"name"`
		Created *TimeStamp            `json:"created"`
		Creator string                `json:"creator"`
	} `json:"channel"`
}

type ChannelDeleted struct {
	CommonEvent
	ChannelID slackobject.ChannelID `json:"channel"`
}

type historyChangedEvent struct {
	Latest         *TimeStamp `json:"latest"`
	TimeStamp      *TimeStamp `json:"ts"`
	EventTimeStamp *TimeStamp `json:"event_ts"`
}

type ChannelHistoryChanged struct {
	CommonEvent
	historyChangedEvent
}

type ChannelJoined struct {
	CommonEvent
	ChannelID slackobject.ChannelID `json:"channel"`
}

type ChannelLeft struct {
	CommonEvent
	ChannelID slackobject.ChannelID `json:"channel"`
}

type markedAsReadEvent struct {
	ChannelID slackobject.ChannelID `json:"channel"`
	TimeStamp *TimeStamp            `json:"ts"`
}

type ChannelMarked struct {
	CommonEvent
	markedAsReadEvent
}

type ChannelRenamed struct {
	CommonEvent
	Channel *struct {
		ID      slackobject.ChannelID `json:"id"`
		Name    string                `json:"name"`
		Created *TimeStamp            `json:"created"`
	} `json:"channel"`
}

type ChannelUnarchived struct {
	CommonEvent
	ChannelID slackobject.ChannelID `json:"channel"`
	User      slackobject.UserID    `json:"user"`
}

type CommandsChanged struct {
	CommonEvent
	TimeStamp *TimeStamp `json:"event_ts"`
}

type DNDUpdated struct {
	CommonEvent
	User      slackobject.UserID `json:"user"`
	DNDStatus struct {
		Enabled            bool       `json:"dnd_enabled"`
		NextStartTimeStamp *TimeStamp `json:"next_dnd_start_ts"`
		NextEndTimeStamp   *TimeStamp `json:"next_dnd_end_ts"`
		SnoozeEnabled      bool       `json:"snooze_enabled"`
		SnoozeEndTimeStamp int        `json:"snooze_endtime"`
	} `json:"dnd_status"`
}

type DNDUpdatedUser struct {
	CommonEvent
	User      slackobject.UserID `json:"user"`
	DNDStatus struct {
		Enabled            bool `json:"dnd_enabled"`
		NextStartTimeStamp int  `json:"next_dnd_start_ts"`
		NextEndTimeStamp   int  `json:"next_dnd_end_ts"`
	} `json:"dnd_status"`
}

type EmailDomainChanged struct {
	CommonEvent
	EmailDomain string     `json:"email_domain"`
	TimeStamp   *TimeStamp `json:"event_ts"`
}

type EmojiChanged struct {
	CommonEvent
	Subtype   string     `json:"subtype"` // TODO add/remove
	Names     []string   `json:"names"`
	TimeStamp *TimeStamp `json:"event_ts"`
}

type Comment struct {
	ID      string             `json:"id"`
	Created *TimeStamp         `json:"created"`
	User    slackobject.UserID `json:"user"`
	Content string             `json:"comment"`
}

type File struct {
	ID                 string                  `json:"id"`
	TimeStamp          *TimeStamp              `json:"created"`
	Name               string                  `json:"name"`
	Title              string                  `json:"title"`
	MimeType           string                  `json:"mimetype"`
	FileType           string                  `json:"filetype"`
	PrettyType         string                  `json:"pretty_type"`
	User               slackobject.UserID      `json:"user"`
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
	Channels           []slackobject.ChannelID `json:"channels"`
	Groups             []string                `json:"groups"`
	Ims                []string                `json:"ims"`
	InitialComment     *Comment                `json:"initial_comment"`
	NumStars           int                     `json:"num_stars"`
	IsStarred          bool                    `json:"is_starred"`
	PinnedTo           []string                `json:"pinned_to"`
	Reactions          []struct {
		Name  string               `json:"name"`
		Count int                  `json:"count"`
		Users []slackobject.UserID `json:"users"`
	} `json:"reactions"`
	CommentsCount int `json:"comments_count"`
}

type FileChanged struct {
	CommonEvent
	File *File `json:"file"`
}

type FileCommentAdded struct {
	CommonEvent
	File    *File    `json:"file"`
	Comment *Comment `json:"comment"`
}

type FileCommentDeleted struct {
	CommonEvent
	File      *File  `json:"file"`
	CommentID string `json:"comment"`
}

type FileCommentEdited struct {
	CommonEvent
	File    *File    `json:"file"`
	Comment *Comment `json:"comment"`
}

type FileCreated struct {
	CommonEvent
	File *File `json:"file"`
}

type FileDeleted struct {
	CommonEvent
	FileID    string     `json:"file_id"`
	TimeStamp *TimeStamp `json:"event_ts"`
}

type FilePublicated struct {
	CommonEvent
	File *File `json:"file"`
}

type FileShared struct {
	CommonEvent
	File *File `json:"file"`
}

type FileUnshared struct {
	CommonEvent
	File *File `json:"file"`
}

type GoodBye struct {
	CommonEvent
}

type GroupArchived struct {
	CommonEvent
	ChannelID slackobject.ChannelID `json:"channel"`
}

type GroupClosed struct {
	CommonEvent
	User      slackobject.UserID    `json:"user"`
	ChannelID slackobject.ChannelID `json:"channel"`
}

type GroupHistoryChanged struct {
	CommonEvent
	historyChangedEvent
}

type GroupJoined struct {
	CommonEvent
	ChannelID slackobject.ChannelID `json:"channel"`
}

type GroupLeft struct {
	CommonEvent
	ChannelID slackobject.ChannelID `json:"channel"`
}

type GroupMarked struct {
	CommonEvent
	markedAsReadEvent
}

type GroupOpened struct {
	CommonEvent
	User      slackobject.UserID    `json:"user"`
	ChannelID slackobject.ChannelID `json:"channel"`
}

type GroupRenamed struct {
	CommonEvent
	Channel struct {
		ID      slackobject.ChannelID `json:"id"`
		Name    string                `json:"name"`
		Created *TimeStamp            `json:"created"`
	} `json:"channel"`
}

type GroupUnarchived struct {
	CommonEvent
	ChannelID slackobject.ChannelID `json:"channel"`
}

// Hello event is sent from slack when WebSocket connection is successfully established.
// https://api.slack.com/events/hello
type Hello struct {
	CommonEvent
}

type IMClosed struct {
	CommonEvent
	User      slackobject.UserID    `json:"user"`
	ChannelID slackobject.ChannelID `json:"channel"`
}

type IMCreated struct {
	CommonEvent
	User    slackobject.UserID `json:"user"`
	Channel struct {
		ID      slackobject.ChannelID `json:"id"`
		Name    string                `json:"name"`
		Created *TimeStamp            `json:"created"`
		Creator string                `json:"creator"`
	} `json:"channel"`
}

type IMHistoryChanged struct {
	CommonEvent
	historyChangedEvent
}

type IMMarked struct {
	CommonEvent
	markedAsReadEvent
}

type IMOpened struct {
	CommonEvent
	User      slackobject.UserID    `json:"user"`
	ChannelID slackobject.ChannelID `json:"channel"`
}

type PresenceManuallyChanged struct {
	CommonEvent
	Presence string `json:"presence"` // TODO Actual values other than "away"
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
	CommonEvent
	ChannelID slackobject.ChannelID `json:"channel"`
	Sender    slackobject.UserID    `json:"user"`
	Text      string                `json:"text"`
	TimeStamp *TimeStamp            `json:"ts"`
}

// Item can be any object with type of Message, File, or Comment.
type Item struct {
	Type      string     `json:"type"`
	Channel   string     `json:"channel"`
	Message   *Message   `json:"message"`
	File      *File      `json:"file"`
	Comment   *Comment   `json:"comment"`
	TimeStamp *TimeStamp `json:"ts"`
}

type PinAdded struct {
	CommonEvent
	User      slackobject.UserID    `json:"user"`
	ChannelID slackobject.ChannelID `json:"channel_id"`
	TimeStamp *TimeStamp            `json:"event_ts"`
	Item      *Item                 `json:"item"`
}

type PinRemoved struct {
	CommonEvent
	User      slackobject.UserID    `json:"user"`
	ChannelID slackobject.ChannelID `json:"channel_id"`
	Item      *Item                 `json:"item"`
	HasPins   bool                  `json:"has_pins"`
	TimeStamp *TimeStamp            `json:"event_ts"`
}

type PreferenceChanged struct {
	CommonEvent
	Name  string `json:"messages_theme"`
	Value string `json:"value"`
}

type PresenceChange struct {
	CommonEvent
	User     slackobject.UserID `json:"user"`
	Presence string             `json:"presence"`
}

type ReactionAdded struct {
	CommonEvent
	User      slackobject.UserID `json:"user"`
	Reaction  string             `json:"reaction"` // TODO actual value
	ItemOwner string             `json:"item_user"`
	Item      *Item              `json:"item"` // TODO message, file, file comment. only ids are given, right?
	TimeStamp *TimeStamp         `json:"event_ts"`
}

type ReactionRemoved struct {
	CommonEvent
	User      slackobject.UserID `json:"user"`
	Reaction  string             `json:"reaction"` // TODO actual value
	ItemOwner string             `json:"item_user"`
	Item      *Item              `json:"item"` // TODO message, file, file comment. only ids are given, right?
	TimeStamp *TimeStamp         `json:"event_ts"`
}

// ReconnectURL is currently unsupported and experimental
// https://api.slack.com/events/reconnect_url
type ReconnectURL struct {
	CommonEvent
}

type StarAdded struct {
	CommonEvent
	User      slackobject.UserID `json:"user"`
	Item      *Item              `json:"item"` // TODO message, file, file comment. only ids are given, right?
	TimeStamp *TimeStamp         `json:"event_ts"`
}

type StarRemoved struct {
	CommonEvent
	User      slackobject.UserID `json:"user"`
	Item      *Item              `json:"item"` // TODO message, file, file comment. only ids are given, right?
	TimeStamp *TimeStamp         `json:"event_ts"`
}

type SubTeam struct {
	ID          string     `json:"id"`
	TeamID      string     `json:"team_id"`
	IsUserGroup bool       `json:"is_usergroup"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Handle      string     `json:"handle"`
	IsExternal  bool       `json:"is_external"`
	Created     *TimeStamp `json:"date_create"`
	Updated     *TimeStamp `json:"date_update"`
	Deleted     *TimeStamp `json:"date_delete"`
	CreatedBy   string     `json:"created_by"`
	UpdatedBy   string     `json:"updated_by"`
	UserCount   int        `json:"user_count"`
}

type SubTeamCreated struct {
	CommonEvent
	SubTeam *SubTeam `json:"subteam"`
}

type SubTeamSelfAdded struct {
	CommonEvent
	SubTeamID string `json:"subteam_id"`
}

type SubTeamSelfRemoved struct {
	CommonEvent
	SubTeamID string `json:"subteam_id"`
}

type SubTeamUpdated struct {
	CommonEvent
	SubTeam *SubTeam `json:"subteam"`
}

type TeamDomainChanged struct {
	CommonEvent
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
	CommonEvent
	User *User `json:"user"`
}

// TeamMigrationStarted is sent when chat group is migrated between servers.
// "The WebSocket connection will close immediately after it is sent.
// *snip* By the time a client has reconnected the process is usually complete, so the impact is minimal."
// https://api.slack.com/events/team_migration_started
type TeamMigrationStarted struct {
	CommonEvent
}

type TeamPlanChanged struct {
	CommonEvent
	Plan string `json:"plan"` // currently "", std, and plus
}

type TeamPreferenceChanged struct {
	CommonEvent
	Name  string `json:"name"`
	Value string `json:"value"`
}

type TeamProfileChanged struct {
	CommonEvent
	Profile *struct {
		Fields []*struct {
			ID string `json:"id"` // TODO
		} `json:"fields"`
	} `json:"profile"`
}

type TeamProfileDeleted struct {
	CommonEvent
	Profile *struct {
		Fields []*struct {
			ID string `json:"id"` // TODO
		} `json:"fields"`
	} `json:"profile"`
}

type TeamProfileReordered struct {
	CommonEvent
	Profile *struct {
		Fields []*struct {
			ID    string `json:"id"` // TODO
			Order int    `json:"ordering"`
		} `json:"fields"`
	} `json:"profile"`
}

type TeamRenamed struct {
	CommonEvent
	Name string `json:"name"`
}

type UserChanged struct {
	CommonEvent
	User *User `json:"user"`
}

type UserTyping struct {
	CommonEvent
	ChannelID slackobject.ChannelID `json:"channel"`
	User      slackobject.UserID    `json:"user"`
}

// Pong is given when client send Ping.
// https://api.slack.com/rtm#ping_and_pong
type Pong struct {
	CommonEvent
	ReplyTo uint `json:"reply_to"`
}

// MiscMessage represents some minor message events.
// TODO define each one with subtype field. This is just a representation of common subtyped payload
// https://api.slack.com/events/message#message_subtypes
type MiscMessage struct {
	CommonMessage
	TimeStamp *TimeStamp `json:"ts"`
}
