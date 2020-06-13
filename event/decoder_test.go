package event

import (
	"github.com/oklahomer/golack/testutil"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

var expectedPayloads = map[string]interface{}{
	"accounts_changed": &AccountsChanged{
		TypedEvent{
			Type: "accounts_changed",
		},
	},
	"bot_added": &BotAdded{
		TypedEvent: TypedEvent{
			Type: "bot_added",
		},
		Bot: &Bot{
			ID:    "B024BE7LH",
			AppID: "A4H1JB4AZ",
			Name:  "hugbot",
			Icon: &struct {
				Image48 string `json:"image_48"`
			}{
				Image48: "https://slack.com/path/to/hugbot_48.png",
			},
		},
	},
	"bot_changed": &BotChanged{
		TypedEvent: TypedEvent{
			Type: "bot_changed",
		},
		Bot: &Bot{
			ID:    "B024BE7LH",
			AppID: "A4H1JB4AZ",
			Name:  "hugbot",
			Icon: &struct {
				Image48 string `json:"image_48"`
			}{
				Image48: "https://slack.com/path/to/hugbot_48.png",
			},
		},
	},
	"channel_archive": &ChannelArchived{
		TypedEvent: TypedEvent{
			Type: "channel_archive",
		},
		ChannelID: "C024BE91L",
		UserID:    "U024BE7LH",
	},
	"channel_created": &ChannelCreated{
		TypedEvent: TypedEvent{
			Type: "channel_created",
		},
		Channel: &struct {
			ID        ChannelID  `json:"id"`
			Name      string     `json:"name"`
			Created   *TimeStamp `json:"created"`
			CreatorID UserID     `json:"creator"`
		}{
			ID:   "C024BE91L",
			Name: "fun",
			Created: &TimeStamp{
				Time:          time.Unix(1360782804, 0),
				OriginalValue: "1360782804",
			},
			CreatorID: "U024BE7LH",
		},
	},
	"channel_deleted": &ChannelDeleted{
		TypedEvent: TypedEvent{
			Type: "channel_deleted",
		},
		ChannelID: "C024BE91L",
	},
	"channel_history_changed": &ChannelHistoryChanged{
		TypedEvent: TypedEvent{
			Type: "channel_history_changed",
		},
		ChangedHistory: ChangedHistory{
			Latest: &TimeStamp{
				Time:          time.Unix(1358877455, 0),
				OriginalValue: "1358877455.000010",
			},
			TimeStamp: &TimeStamp{
				Time:          time.Unix(1361482916, 0),
				OriginalValue: "1361482916.000003",
			},
			EventTimeStamp: &TimeStamp{
				Time:          time.Unix(1361482916, 0),
				OriginalValue: "1361482916.000004",
			},
		},
	},
	"channel_joined": &ChannelJoined{
		TypedEvent: TypedEvent{
			Type: "channel_joined",
		},
		ChannelID: "C024BE91L",
	},
	"channel_left": &ChannelLeft{
		TypedEvent: TypedEvent{
			Type: "channel_left",
		},
		ChannelID: "C024BE91L",
	},
	"channel_marked": &ChannelMarked{
		TypedEvent: TypedEvent{
			Type: "channel_marked",
		},
		MarkedAsRead: MarkedAsRead{
			ChannelID: "C024BE91L",
			TimeStamp: &TimeStamp{
				Time:          time.Unix(1401383885, 0),
				OriginalValue: "1401383885.000061",
			},
		},
	},
	"channel_rename": &ChannelRenamed{
		TypedEvent: TypedEvent{
			Type: "channel_rename",
		},
		Channel: &struct {
			ID      ChannelID  `json:"id"`
			Name    string     `json:"name"`
			Created *TimeStamp `json:"created"`
		}{
			ID:   "C02ELGNBH",
			Name: "new_name",
			Created: &TimeStamp{
				Time:          time.Unix(1360782804, 0),
				OriginalValue: "1360782804",
			},
		},
	},
	"channel_unarchive": &ChannelUnarchived{
		TypedEvent: TypedEvent{
			Type: "channel_unarchive",
		},
		ChannelID: "C024BE91L",
		UserID:    "U024BE7LH",
	},
	"commands_changed": &CommandsChanged{
		TypedEvent: TypedEvent{
			Type: "commands_changed",
		},
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1361482916, 0),
			OriginalValue: "1361482916.000004",
		},
	},
	"dnd_updated": &DNDUpdated{
		TypedEvent: TypedEvent{
			Type: "dnd_updated",
		},
		UserID: "U1234",
		DNDStatus: &struct {
			Enabled            bool       `json:"dnd_enabled"`
			NextStartTimeStamp *TimeStamp `json:"next_dnd_start_ts"`
			NextEndTimeStamp   *TimeStamp `json:"next_dnd_end_ts"`
			SnoozeEnabled      bool       `json:"snooze_enabled"`
			SnoozeEndTimeStamp *TimeStamp `json:"snooze_endtime"`
		}{
			Enabled: true,
			NextStartTimeStamp: &TimeStamp{
				Time:          time.Unix(1450387800, 0),
				OriginalValue: "1450387800",
			},
			NextEndTimeStamp: &TimeStamp{
				Time:          time.Unix(1450423800, 0),
				OriginalValue: "1450423800",
			},
			SnoozeEnabled: true,
			SnoozeEndTimeStamp: &TimeStamp{
				Time:          time.Unix(1450373897, 0),
				OriginalValue: "1450373897",
			},
		},
	},
	"dnd_updated_user": &DNDUpdatedUser{
		TypedEvent: TypedEvent{
			Type: "dnd_updated_user",
		},
		UserID: "U1234",
		DNDStatus: &struct {
			Enabled            bool       `json:"dnd_enabled"`
			NextStartTimeStamp *TimeStamp `json:"next_dnd_start_ts"`
			NextEndTimeStamp   *TimeStamp `json:"next_dnd_end_ts"`
		}{
			Enabled: true,
			NextStartTimeStamp: &TimeStamp{
				Time:          time.Unix(1450387800, 0),
				OriginalValue: "1450387800",
			},
			NextEndTimeStamp: &TimeStamp{
				Time:          time.Unix(1450423800, 0),
				OriginalValue: "1450423800",
			},
		},
	},
	"email_domain_changed": &EmailDomainChanged{
		TypedEvent: TypedEvent{
			Type: "email_domain_changed",
		},
		EmailDomain: "example.com",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1360782804, 0),
			OriginalValue: "1360782804.083113",
		},
	},
	"emoji_changed": &EmojiChanged{
		TypedEvent: TypedEvent{
			Type: "emoji_changed",
		},
		Subtype: "remove",
		Names:   []string{"picard_facepalm"},
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1361482916, 0),
			OriginalValue: "1361482916.000004",
		},
	},
	"file_change": &FileChanged{
		TypedEvent: TypedEvent{
			Type: "file_change",
		},
		FileID: "F2147483862",
		File: &struct {
			ID FileID `json:"id"`
		}{
			ID: "F2147483862",
		},
	},
	"file_comment_added": &FileCommentAdded{
		TypedEvent: TypedEvent{
			Type: "file_comment_added",
		},
		FileID: "F2147483862",
		File: &struct {
			ID FileID `json:"id"`
		}{
			ID: "F2147483862",
		},
		Comment: &Comment{
			ID: "C12345",
			Created: &TimeStamp{
				Time:          time.Unix(1360782804, 0),
				OriginalValue: "1360782804",
			},
			UserID:  "U1234",
			Content: "comment content",
		},
	},
	"file_comment_deleted": &FileCommentDeleted{
		TypedEvent: TypedEvent{
			Type: "file_comment_deleted",
		},
		CommentID: "Fc67890",
		FileID:    "F2147483862",
		File: &struct {
			ID FileID `json:"id"`
		}{
			ID: "F2147483862",
		},
	},
	"file_comment_edited": &FileCommentEdited{
		TypedEvent: TypedEvent{
			Type: "file_comment_edited",
		},
		FileID: "F2147483862",
		File: &struct {
			ID FileID `json:"id"`
		}{
			ID: "F2147483862",
		},
		Comment: &Comment{
			ID: "C12345",
			Created: &TimeStamp{
				Time:          time.Unix(1360782804, 0),
				OriginalValue: "1360782804",
			},
			UserID:  "U1234",
			Content: "comment content",
		},
	},
	"file_created": &FileCreated{
		TypedEvent: TypedEvent{
			Type: "file_created",
		},
		FileID: "F2147483862",
		File: &struct {
			ID FileID `json:"id"`
		}{
			ID: "F2147483862",
		},
	},
	"file_deleted": &FileDeleted{
		TypedEvent: TypedEvent{
			Type: "file_deleted",
		},
		FileID: "F2147483862",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1361482916, 0),
			OriginalValue: "1361482916.000004",
		},
	},
	"file_public": &FilePublished{
		TypedEvent: TypedEvent{
			Type: "file_public",
		},
		FileID: "F2147483862",
		File: &struct {
			ID FileID `json:"id"`
		}{
			ID: "F2147483862",
		},
	},
	"file_shared": &FileShared{
		TypedEvent: TypedEvent{
			Type: "file_shared",
		},
		FileID: "F2147483862",
		File: &struct {
			ID FileID `json:"id"`
		}{
			ID: "F2147483862",
		},
	},
	"file_unshared": &FileUnshared{
		TypedEvent: TypedEvent{
			Type: "file_unshared",
		},
		FileID: "F2147483862",
		File: &struct {
			ID FileID `json:"id"`
		}{
			ID: "F2147483862",
		},
	},
	"goodbye": &GoodBye{
		TypedEvent: TypedEvent{
			Type: "goodbye",
		},
	},
	"group_archive": &GroupArchived{
		TypedEvent: TypedEvent{
			Type: "group_archive",
		},
		ChannelID: "G024BE91L",
	},
	"group_close": &GroupClosed{
		TypedEvent: TypedEvent{
			Type: "group_close",
		},
		UserID:    "U024BE7LH",
		ChannelID: "G024BE91L",
	},
	"group_deleted": &GroupDeleted{
		TypedEvent: TypedEvent{
			Type: "group_deleted",
		},
		ChannelID: "G0QN9RGTT",
	},
	"group_history_changed": &GroupHistoryChanged{
		TypedEvent: TypedEvent{
			Type: "group_history_changed",
		},
		ChangedHistory: ChangedHistory{
			Latest: &TimeStamp{
				Time:          time.Unix(1358877455, 0),
				OriginalValue: "1358877455.000010",
			},
			TimeStamp: &TimeStamp{
				Time:          time.Unix(1361482916, 0),
				OriginalValue: "1361482916.000003",
			},
			EventTimeStamp: &TimeStamp{
				Time:          time.Unix(1361482916, 0),
				OriginalValue: "1361482916.000004",
			},
		},
	},
	"group_joined": &GroupJoined{
		TypedEvent: TypedEvent{
			Type: "group_joined",
		},
		ChannelID: "G0QN9RGTT",
	},
	"group_left": &GroupLeft{
		TypedEvent: TypedEvent{
			Type: "group_left",
		},
		ChannelID: "G02ELGNBH",
	},
	"group_marked": &GroupMarked{
		TypedEvent: TypedEvent{
			Type: "group_marked",
		},
		MarkedAsRead: MarkedAsRead{
			ChannelID: "G024BE91L",
			TimeStamp: &TimeStamp{
				Time:          time.Unix(1401383885, 0),
				OriginalValue: "1401383885.000061",
			},
		},
	},
	"group_open": &GroupOpened{
		TypedEvent: TypedEvent{
			Type: "group_open",
		},
		UserID:    "U024BE7LH",
		ChannelID: "G024BE91L",
	},
	"group_rename": &GroupRenamed{
		TypedEvent: TypedEvent{
			Type: "group_rename",
		},
		Channel: &struct {
			ID      ChannelID  `json:"id"`
			Name    string     `json:"name"`
			Created *TimeStamp `json:"created"`
		}{
			ID:   "G02ELGNBH",
			Name: "new_name",
			Created: &TimeStamp{
				Time:          time.Unix(1360782804, 0),
				OriginalValue: "1360782804",
			},
		},
	},
	"group_unarchive": &GroupUnarchived{
		TypedEvent: TypedEvent{
			Type: "group_unarchive",
		},
		ChannelID: "G024BE91L",
	},
	"hello": &Hello{
		TypedEvent: TypedEvent{
			Type: "hello",
		},
	},
	"im_close": &IMClosed{
		TypedEvent: TypedEvent{
			Type: "im_close",
		},
		UserID:    "U024BE7LH",
		ChannelID: "D024BE91L",
	},
	"im_created": &IMCreated{
		TypedEvent: TypedEvent{
			Type: "im_created",
		},
		UserID: "U024BE7LH",
		Channel: &struct {
			ID ChannelID `json:"id"`
		}{
			ID: "D024BE91L",
		},
	},
	"im_history_changed": &IMHistoryChanged{
		TypedEvent: TypedEvent{
			Type: "im_history_changed",
		},
		ChangedHistory: ChangedHistory{
			Latest: &TimeStamp{
				Time:          time.Unix(1358877455, 0),
				OriginalValue: "1358877455.000010",
			},
			TimeStamp: &TimeStamp{
				Time:          time.Unix(1361482916, 0),
				OriginalValue: "1361482916.000003",
			},
			EventTimeStamp: &TimeStamp{
				Time:          time.Unix(1361482916, 0),
				OriginalValue: "1361482916.000004",
			},
		},
	},
	"im_marked": &IMMarked{
		TypedEvent: TypedEvent{
			Type: "im_marked",
		},
		MarkedAsRead: MarkedAsRead{
			ChannelID: "D024BE91L",
			TimeStamp: &TimeStamp{
				Time:          time.Unix(1401383885, 0),
				OriginalValue: "1401383885.000061",
			},
		},
	},
	"im_open": &IMOpened{
		TypedEvent: TypedEvent{
			Type: "im_open",
		},
		UserID:    "U024BE7LH",
		ChannelID: "D024BE91L",
	},
	"manual_presence_change": &PresenceManuallyChanged{
		TypedEvent: TypedEvent{
			Type: "manual_presence_change",
		},
		Presence: "away",
	},
	"member_joined_channel": &MemberJoinedChannel{
		TypedEvent: TypedEvent{
			Type: "member_joined_channel",
		},
		UserID:      "W06GH7XHN",
		ChannelID:   "C0698JE0H",
		ChannelType: "C",
		TeamID:      "T024BE7LD",
		InviterID:   "U123456789",
	},
	"message": &Message{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		ChannelID: "C2147483705",
		SenderID:  "U2147483697",
		Text:      "Hello world",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1355517523, 0),
			OriginalValue: "1355517523.000005",
		},
	},
	"pin_added": &PinAdded{
		TypedEvent: TypedEvent{
			Type: "pin_added",
		},
		UserID:    "U024BE7LH",
		ChannelID: "C02ELGNBH",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1360782804, 0),
			OriginalValue: "1360782804.083113",
		},
	},
	"pin_removed": &PinRemoved{
		TypedEvent: TypedEvent{
			Type: "pin_removed",
		},
		UserID:    "U024BE7LH",
		ChannelID: "C02ELGNBH",
		HasPins:   false,
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1360782804, 0),
			OriginalValue: "1360782804.083113",
		},
	},
	"pref_change": &PreferenceChanged{
		TypedEvent: TypedEvent{
			Type: "pref_change",
		},
		Name:  "messages_theme",
		Value: "dense",
	},
	"presence_change": &PresenceChanged{
		TypedEvent: TypedEvent{
			Type: "presence_change",
		},
		UserID:   "U024BE7LH",
		Presence: "away",
	},
	"presence_query": &PresenceQuery{
		TypedEvent: TypedEvent{
			Type: "presence_query",
		},
		UserIDs: []UserID{
			"U061F7AUR",
			"W123456",
		},
	},
	"presence_sub": &PresenceSubscribe{
		TypedEvent: TypedEvent{
			Type: "presence_sub",
		},
		UserIDs: []UserID{
			"U061F7AUR",
			"W123456",
		},
	},
	"reaction_added": &ReactionAdded{
		TypedEvent: TypedEvent{
			Type: "reaction_added",
		},
		UserID:      "U024BE7LH",
		Reaction:    "thumbsup",
		ItemOwnerID: "U0G9QF9C6",
		Item: &Item{
			Type:      "message",
			ChannelID: "C0G9QF9GZ",
			TimeStamp: &TimeStamp{
				Time:          time.Unix(1360782400, 0),
				OriginalValue: "1360782400.498405",
			},
		},
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1360782804, 0),
			OriginalValue: "1360782804.083113",
		},
	},
	"reaction_removed": &ReactionRemoved{
		TypedEvent: TypedEvent{
			Type: "reaction_removed",
		},
		UserID:      "U024BE7LH",
		Reaction:    "thumbsup",
		ItemOwnerID: "U0G9QF9C6",
		Item: &Item{
			Type:      "message",
			ChannelID: "C0G9QF9GZ",
			TimeStamp: &TimeStamp{
				Time:          time.Unix(1360782400, 0),
				OriginalValue: "1360782400.498405",
			},
		},
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1360782804, 0),
			OriginalValue: "1360782804.083113",
		},
	},
	"reconnect_url": &ReconnectURL{
		TypedEvent: TypedEvent{
			Type: "reconnect_url",
		},
	},
	"star_added": &StarAdded{
		TypedEvent: TypedEvent{
			Type: "star_added",
		},
		UserID: "U024BE7LH",
		Item: &Item{
			Type:      "message",
			ChannelID: "C0G9QF9GZ",
			TimeStamp: &TimeStamp{
				Time:          time.Unix(1360782400, 0),
				OriginalValue: "1360782400.498405",
			},
		},
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1360782804, 0),
			OriginalValue: "1360782804.083113",
		},
	},
	"star_removed": &StarRemoved{
		TypedEvent: TypedEvent{
			Type: "star_removed",
		},
		UserID: "U024BE7LH",
		Item: &Item{
			Type:      "message",
			ChannelID: "C0G9QF9GZ",
			TimeStamp: &TimeStamp{
				Time:          time.Unix(1360782400, 0),
				OriginalValue: "1360782400.498405",
			},
		},
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1360782804, 0),
			OriginalValue: "1360782804.083113",
		},
	},
	"subteam_created": &SubTeamCreated{
		TypedEvent: TypedEvent{
			Type: "subteam_created",
		},
		SubTeam: &SubTeam{
			ID:          "S0615G0KT",
			TeamID:      "T060RNRCH",
			IsUserGroup: true,
			Name:        "Marketing Team",
			Description: "Marketing gurus, PR experts and product advocates.",
			Handle:      "marketing-team",
			IsExternal:  false,
			Created: &TimeStamp{
				Time:          time.Unix(1446746793, 0),
				OriginalValue: "1446746793",
			},
			Updated: &TimeStamp{
				Time:          time.Unix(1446746793, 0),
				OriginalValue: "1446746793",
			},
			Deleted: &TimeStamp{
				Time:          time.Unix(1446746793, 0),
				OriginalValue: "1446746793",
			},
			AutoType:  "",
			CreatorID: "U060RNRCZ",
			UpdatedBy: "U060RNRCZ",
			UserCount: 10,
		},
	},
	"subteam_members_changed": &SubTeamMembersChanged{
		TypedEvent: TypedEvent{
			Type: "subteam_members_changed",
		},
		SubTeamID: "S0614TZR7",
		TeamID:    "T060RNRCH",
		PreviousUpdate: &TimeStamp{
			Time:          time.Unix(1446670362, 0),
			OriginalValue: "1446670362",
		},
		Updated: &TimeStamp{
			Time:          time.Unix(1492906952, 0),
			OriginalValue: "1492906952",
		},
		AddedUserIDs: []UserID{
			"U060RNRCZ",
			"U060ULRC0",
			"U061309JM",
		},
		AddedUserCount: 3,
		RemovedUserIDs: []UserID{
			"U06129G2V",
		},
		RemovedUsersCount: 1,
	},
	"subteam_self_added": &SubTeamSelfAdded{
		TypedEvent: TypedEvent{
			Type: "subteam_self_added",
		},
		SubTeamID: "S0615G0KT",
	},
	"subteam_self_removed": &SubTeamSelfRemoved{
		TypedEvent: TypedEvent{
			Type: "subteam_self_removed",
		},
		SubTeamID: "S0615G0KT",
	},
	"subteam_updated": &SubTeamUpdated{
		TypedEvent: TypedEvent{
			Type: "subteam_updated",
		},
		SubTeam: &SubTeam{
			ID:          "S0614TZR7",
			TeamID:      "T060RNRCH",
			IsUserGroup: true,
			Name:        "Team Admins",
			Description: "A group of all Administrators on your team.",
			Handle:      "admins",
			IsExternal:  false,
			Created: &TimeStamp{
				Time:          time.Unix(1446598059, 0),
				OriginalValue: "1446598059",
			},
			Updated: &TimeStamp{
				Time:          time.Unix(1446670362, 0),
				OriginalValue: "1446670362",
			},
			Deleted: &TimeStamp{
				Time:          time.Unix(0, 0),
				OriginalValue: "0",
			},
			AutoType:  "admin",
			CreatorID: "USLACKBOT",
			UpdatedBy: "U060RNRCZ",
			UserCount: 4,
			UserIDs: []UserID{
				"U060RNRCZ",
				"U060ULRC0",
				"U06129G2V",
				"U061309JM",
			},
		},
	},
	"team_domain_change": &TeamDomainChanged{
		TypedEvent: TypedEvent{
			Type: "team_domain_change",
		},
		URL:    "https://my.slack.com",
		Domain: "my",
	},
	"team_join": &TeamJoined{
		TypedEvent: TypedEvent{
			Type: "team_join",
		},
		User: &User{},
	},
	"team_migration_started": &TeamMigrationStarted{
		TypedEvent: TypedEvent{
			Type: "team_migration_started",
		},
	},
	"team_plan_change": &TeamPlanChanged{
		TypedEvent: TypedEvent{
			Type: "team_plan_change",
		},
		Plan:      "std",
		CanAddUra: false,
		PaidFeatures: []string{
			"feature1",
			"feature2",
		},
	},
	"team_pref_change": &TeamPreferenceChanged{
		TypedEvent: TypedEvent{
			Type: "team_pref_change",
		},
		Name:  "slackbot_responses_only_admins",
		Value: true,
	},
	"team_profile_change": &TeamProfileChanged{
		TypedEvent: TypedEvent{
			Type: "team_profile_change",
		},
		Profile: &struct {
			Fields []*struct {
				ID string `json:"id"`
			} `json:"fields"`
		}{
			Fields: []*struct {
				ID string `json:"id"`
			}{
				{
					ID: "Xf06054AAA",
				},
			},
		},
	},
	"team_profile_delete": &TeamProfileDeleted{
		TypedEvent: TypedEvent{
			Type: "team_profile_delete",
		},
		Profile: &struct {
			Fields []string `json:"fields"`
		}{
			Fields: []string{
				"Xf06054AAA",
			},
		},
	},
	"team_profile_reorder": &TeamProfileReordered{
		TypedEvent: TypedEvent{
			Type: "team_profile_reorder",
		},
		Profile: &struct {
			Fields []*struct {
				ID    string `json:"id"`
				Order int    `json:"ordering"`
			} `json:"fields"`
		}{
			Fields: []*struct {
				ID    string `json:"id"`
				Order int    `json:"ordering"`
			}{
				{
					ID:    "Xf06054AAA",
					Order: 0,
				},
			},
		},
	},
	"team_rename": &TeamRenamed{
		TypedEvent: TypedEvent{
			Type: "team_rename",
		},
		Name: "New Team Name Inc.",
	},
	"user_change": &UserChanged{
		TypedEvent: TypedEvent{
			Type: "user_change",
		},
		User: &User{},
	},
	"user_typing": &UserTyping{
		TypedEvent: TypedEvent{
			Type: "user_typing",
		},
		ChannelID: "C02ELGNBH",
		UserID:    "U024BE7LH",
	},
}

func TestDecode(t *testing.T) {
	directory := filepath.Join("..", "testdata", "event", "decode")

	checked := map[string]bool{}
	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Fatalf("Failed to read testdata directory: %s", err.Error())
		}

		if !strings.HasSuffix(path, ".json.golden") {
			// Skip irrelevant file
			return nil
		}

		filename := filepath.Base(path)
		t.Run(filename, func(t *testing.T) {
			input, err := ioutil.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to read file: %s. Error: %s.", path, err.Error())
			}

			eventName := strings.TrimSuffix(filename, ".json.golden")
			checked[eventName] = true

			decoded, err := Decode(input)
			if err != nil {
				t.Fatalf("Failed to decode JSON: %s. Error: %s.", path, err.Error())
			}

			typed, ok := decoded.(Typer)
			if !ok {
				t.Errorf("Decoded paylaod is not typed event: %#v", input)
				return
			}

			expected, ok := expectedPayloads[typed.EventType()]
			if !ok {
				t.Fatalf("Expected payload for %s is not defined in test", typed.EventType())
			}

			testutil.Compare([]string{typed.EventType()}, reflect.ValueOf(expected), reflect.ValueOf(typed), t)
		})

		return nil
	})

	var missing []string
	for k := range eventTypeMap {
		_, ok := checked[k]
		if !ok {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		t.Errorf("Tests for %s are missing", missing)
	}
	// TODO add missing test for subtypes

	t.Run("empty string", func(t *testing.T) {
		_, err := Decode([]byte(` `))
		if err != ErrEmptyPayload {
			t.Fatalf("Expected error is not returned: %#v", err)
		}
	})

	t.Run("unknown type", func(t *testing.T) {
		_, err := Decode([]byte(`{"type": "UNKNOWN_VALUE"`))
		if err == nil {
			t.Fatalf("Expected error is not returned: %#v", err)
		}

		typedErr, ok := err.(*UnknownPayloadTypeError)
		if !ok {
			t.Fatalf("Returned error is not UnknownPayloadTypeError but %T", err)
		}

		if !strings.Contains(typedErr.Error(), "UNKNOWN_VALUE") {
			t.Errorf("Given error string does not contain field name: %s", typedErr.Error())
		}
	})

	t.Run("unknown subtype", func(t *testing.T) {
		_, err := Decode([]byte(`{"type": "message", "subtype": "UNKNOWN_VALUE"}`))
		if err == nil {
			t.Fatalf("Expected error is not returned.")
		}

		typedErr, ok := err.(*UnknownPayloadTypeError)
		if !ok {
			t.Fatalf("Returned error is not UnknownPayloadTypeError but %T", err)
		}

		if !strings.Contains(typedErr.Error(), "UNKNOWN_VALUE") {
			t.Errorf("Given error string does not contain field name: %s", typedErr.Error())
		}
	})
}
