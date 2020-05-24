package event

// TypedEvent takes care of events that have "type" field in its JSON representation.
//
// ARM API document, https://api.slack.com/rtm#events, states as follows:
// "Every event has a type property which describes the type of event."
//
// Similarly, Events API document, https://api.slack.com/events-api#event_type_structure, states as follows:
// "The specific name of the event described by its adjacent fields. This field is included with every inner event type."
type TypedEvent struct {
	Type string `json:"type,omitempty"`
}

func (t TypedEvent) EventType() string {
	return t.Type
}

type Typer interface {
	EventType() string
}

// https://api.slack.com/events

// AccountsChanged event indicates that the list of accounts a user is signed into has changed.
// https://api.slack.com/events/accounts_changed
type AccountsChanged struct {
	TypedEvent
}

// TODO add app_home_opened
// TODO add app_mention
// TODO add app_rate_limited
// TODO app_requested
// TODO app_uninstalled

type BotAdded struct {
	TypedEvent
	Bot *Bot `json:"bot"`
}

type BotChanged struct {
	TypedEvent
	Bot *Bot `json:"bot"`
}

// TODO call_rejected

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

// TODO channel_shared

type ChannelRenamed struct {
	TypedEvent
	Channel *struct {
		ID      ChannelID  `json:"id"`
		Name    string     `json:"name"`
		Created *TimeStamp `json:"created"`
	} `json:"channel"`
}

type ChannelUnarchived struct {
	TypedEvent
	ChannelID ChannelID `json:"channel"`
	UserID    UserID    `json:"user"`
}

// TODO channel_unshared

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

// TODO external_org_migration_finished
// TODO external_org_migration_started

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

// TODO add grid_migration_finished
// TODO grid_migration_started

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

// TODO add invite_requested
// TODO add link_shared

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

// TODO add member_left_channel

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
	ChannelID       ChannelID  `json:"channel"`
	SenderID        UserID     `json:"user"`
	Text            string     `json:"text"`
	TimeStamp       *TimeStamp `json:"ts"`
	ThreadTimeStamp *TimeStamp `json:"thread_ts"` // https://api.slack.com/docs/message-threading
}

// TODO add message.app_home
// TODO add message.channels
// TODO add message.groups
// TODO add message.im
// TODO add message.mpim

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

// TODO add resources_added
// TODO add resources_removed
// TODO add scope_denied
// TODO add scope_granted

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

// TODO add tokens_revoked
// TODO add url_verification

type UserChanged struct {
	TypedEvent
	User *User `json:"user"`
}

// TODO add user_resource_denied
// TODO add user_resource_granted
// TODO add user_resource_removed

type UserTyping struct {
	TypedEvent
	ChannelID ChannelID `json:"channel"`
	UserID    UserID    `json:"user"`
}

// Below comes some types that are commonly shared amongst multiple events

type Bot struct {
	ID    BotID  `json:"id"`
	AppID AppID  `json:"app_id"`
	Name  string `json:"name"`
	Icon  *struct {
		Image48 string `json:"image_48"`
	} `json:"icons"`
}

type ChangedHistory struct {
	Latest         *TimeStamp `json:"latest"`
	TimeStamp      *TimeStamp `json:"ts"`
	EventTimeStamp *TimeStamp `json:"event_ts"`
}

type Comment struct {
	ID      CommentID  `json:"id"`
	Created *TimeStamp `json:"created"`
	UserID  UserID     `json:"user"`
	Content string     `json:"comment"`
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
