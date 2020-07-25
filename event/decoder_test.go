package event

import (
	"github.com/oklahomer/golack/v2/testutil"
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
	"app_home_opened": &AppHomeOpened{
		TypedEvent: TypedEvent{
			Type: "app_home_opened",
		},
		UserID:    "U061F7AUR",
		ChannelID: "D0LAN2Q65",
		EventTimeStamp: &TimeStamp{
			Time:          time.Unix(1515449522000016, 0),
			OriginalValue: "1515449522000016",
		},
		Tab: "home",
		View: &View{
			ID:     "VPASKP233",
			TeamID: "T21312902",
			Type:   "home",
			Blocks: []Block{
				&InputBlock{
					block: block{
						Type:    "input",
						BlockID: "multi-line",
					},
					Label: &TextCompositionObject{
						Type: "plain_text",
						Text: "Enter your value",
					},
					Element: &PlainTextInputBlockElement{
						blockElement: blockElement{
							Type: "plain_text_input",
						},
						Multiline: true,
						ActionID:  "ml-value",
					},
					Hint:     nil,
					Optional: false,
				},
			},
			PrivateMetadata: "",
			CallbackID:      "",
			State: &ViewState{
				Values: map[BlockID]map[ActionID]*ViewStateValue{
					"multi-line": {
						"ml-value": &ViewStateValue{
							Type:  "plain_text_input",
							Value: "This is my example inputted value",
						},
					},
				},
			},
			Hash:               "1231232323.12321312",
			ClearOnClose:       false,
			NotifyOnClose:      false,
			RootViewID:         "VPASKP233",
			AppID:              "A21SDS90",
			ExternalID:         "",
			AppInstalledTeamID: "T21312902",
			BotID:              "BSDKSAO2",
		},
	},
	"app_mention": &AppMention{
		TypedEvent: TypedEvent{
			Type: "app_mention",
		},
		UserID: "U061F7AUR",
		Text:   "<@U0LAN0Z89> is it everything a river should be?",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1515449522, 0),
			OriginalValue: "1515449522.000016",
		},
		ChannelID: "C0LAN2Q65",
		EventTimeStamp: &TimeStamp{
			Time:          time.Unix(1515449522000016, 0),
			OriginalValue: "1515449522000016",
		},
	},
	"app_rate_limited": &AppRateLimited{
		TypedEvent: TypedEvent{
			Type: "app_rate_limited",
		},
		Token:  "Jhj5dZrVaK7ZwHHjRyZWjbDl",
		TeamID: "T123456",
		MinuteRateLimited: &TimeStamp{
			Time:          time.Unix(1518467820, 0),
			OriginalValue: "1518467820",
		},
		APIAppID: "A123456",
	},
	"app_requested": &AppRequested{
		TypedEvent: TypedEvent{
			Type: "app_requested",
		},
		AppRequest: &AppRequest{
			ID: "1234",
			App: &App{
				ID:                     "A5678",
				Name:                   "Brent's app",
				Description:            "They're good apps, Bront.",
				HelpURL:                "brontsapp.com",
				PrivacyPolicyURL:       "brontsapp.com",
				AppHomepageURL:         "brontsapp.com",
				AppDirectoryURL:        "https://slack.slack.com/apps/A102ARD7Y",
				IsAppDirectoryApproved: true,
				IsInternal:             false,
				AdditionalInfo:         "none",
			},
			PreviousResolution: &struct {
				Status string             `json:"status"`
				Scopes []*AppRequestScope `json:"scopes"`
			}{
				Status: "approved",
				Scopes: []*AppRequestScope{
					{
						Name:        "app_requested",
						Description: "allows this app to listen for app install requests",
						IsSensitive: false,
						TokenType:   "user",
					},
				},
			},
			User: &struct {
				ID    UserID `json:"id"`
				Name  string `json:"name"`
				Email string `json:"email"`
			}{
				ID:    "U1234",
				Name:  "Bront",
				Email: "bront@brent.com",
			},
			Team: &struct {
				ID     TeamID `json:"id"`
				Name   string `json:"name"`
				Domain string `json:"domain"`
			}{
				ID:     "T1234",
				Name:   "Brant App Team",
				Domain: "brantappteam",
			},
			Scopes: []*AppRequestScope{
				{
					Name:        "app_requested",
					Description: "allows this app to listen for app install requests",
					IsSensitive: false,
					TokenType:   "user",
				},
			},
			Message: "none",
		},
	},
	"app_uninstalled": &AppUninstalled{
		TypedEvent: TypedEvent{
			Type: "app_uninstalled",
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
			Icon: &BotIcon{
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
			Icon: &BotIcon{
				Image48: "https://slack.com/path/to/hugbot_48.png",
			},
		},
	},
	"call_rejected": &CallRejected{
		TypedEvent: TypedEvent{
			Type: "call_rejected",
		},
		CallID:           "RL731AVEF",
		UserID:           "ULJS1TYR5",
		ChannelID:        "DL5JN9K0T",
		ExternalUniqueID: "123-456-7890",
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
	"channel_shared": &ChannelShared{
		TypedEvent: TypedEvent{
			Type: "channel_shared",
		},
		ConnectedTeamID: "E163Q94DX",
		ChannelID:       "C06EQBRR6",
		EventTimeStamp: &TimeStamp{
			Time:          time.Unix(1561064063, 0),
			OriginalValue: "1561064063.001100",
		},
	},
	"channel_unarchive": &ChannelUnarchived{
		TypedEvent: TypedEvent{
			Type: "channel_unarchive",
		},
		ChannelID: "C024BE91L",
		UserID:    "U024BE7LH",
	},
	"channel_unshared": &ChannelUnshared{
		TypedEvent: TypedEvent{
			Type: "channel_unshared",
		},
		PreviouslyConnectedTeamID: "E163Q94DX",
		ChannelID:                 "C06EQBRR6",
		IsExtShared:               false,
		EventTimeStamp: &TimeStamp{
			Time:          time.Unix(1561064063, 0),
			OriginalValue: "1561064063.001100",
		},
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
	"external_org_migration_finished": &ExternalOrgMigrationFinished{
		TypedEvent: TypedEvent{
			Type: "external_org_migration_finished",
		},
		Team: &struct {
			ID          TeamID `json:"id"`
			IsMigrating bool   `json:"is_migrating"`
		}{
			ID:          "TXXXXXXXX",
			IsMigrating: false,
		},
		DateStarted: &TimeStamp{
			Time:          time.Unix(1551398400, 0),
			OriginalValue: "1551398400",
		},
		DateFinished: &TimeStamp{
			Time:          time.Unix(1551409200, 0),
			OriginalValue: "1551409200",
		},
	},
	"external_org_migration_started": &ExternalOrgMigrationStarted{
		TypedEvent: TypedEvent{
			Type: "external_org_migration_started",
		},
		Team: &struct {
			ID          TeamID `json:"id"`
			IsMigrating bool   `json:"is_migrating"`
		}{
			ID:          "TXXXXXXXX",
			IsMigrating: true,
		},
		DateStarted: &TimeStamp{
			Time:          time.Unix(1551398400, 0),
			OriginalValue: "1551398400",
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
			TimeStamp: &TimeStamp{
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
			TimeStamp: &TimeStamp{
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
	"grid_migration_finished": &GridMigrationFinished{
		TypedEvent: TypedEvent{
			Type: "grid_migration_finished",
		},
		EnterpriseID: "EXXXXXXXX",
	},
	"grid_migration_started": &GridMigrationStarted{
		TypedEvent: TypedEvent{
			Type: "grid_migration_started",
		},
		EnterpriseID: "EXXXXXXXX",
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
	"invite_requested": &InviteRequested{
		TypedEvent: TypedEvent{
			Type: "invite_requested",
		},
		InviteRequest: &struct {
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
		}{
			ID:    "12345",
			Email: "bront@puppies.com",
			DateCreated: &TimeStamp{
				Time:          time.Unix(123455, 0),
				OriginalValue: "123455",
			},
			RequesterIDs: []UserID{"U12345"},
			ChannelIDs:   []ChannelID{"C12345"},
			InviteType:   "full_member",
			RealName:     "Brent",
			DateExpire: &TimeStamp{
				Time:          time.Unix(123456, 0),
				OriginalValue: "123456",
			},
			RequestReason: "They're good dogs, Brant",
			Team: &struct {
				ID     TeamID `json:"id"`
				Name   string `json:"name"`
				Domain string `json:"domain"`
			}{
				ID:     "T12345",
				Name:   "Puppy ratings workspace incorporated",
				Domain: "puppiesrus",
			},
		},
	},
	"link_shared": &LinkShared{
		TypedEvent: TypedEvent{
			Type: "link_shared",
		},
		ChannelID: "Cxxxxxx",
		UserID:    "Uxxxxxxx",
		MessageTimeStamp: &TimeStamp{
			Time:          time.Unix(123456789, 0),
			OriginalValue: "123456789.9875",
		},
		ThreadTimeStamp: &TimeStamp{
			Time:          time.Unix(123456621, 0),
			OriginalValue: "123456621.1855",
		},
		Links: []*struct {
			Domain string `json:"domain"`
			URL    string `json:"url"`
		}{
			{
				Domain: "example.com",
				URL:    "https://example.com/12345",
			},
			{
				Domain: "example.com",
				URL:    "https://example.com/67890",
			},
			{
				Domain: "another-example.com",
				URL:    "https://yet.another-example.com/v/abcde",
			},
		},
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
	"member_left_channel": &MemberLeftChannel{
		TypedEvent: TypedEvent{
			Type: "member_left_channel",
		},
		UserID:      "W06GH7XHN",
		ChannelID:   "C0698JE0H",
		ChannelType: "C",
		TeamID:      "T024BE7LD",
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
	"message.app_home": &MessageAppHome{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		UserID: "U061F7AUR",
		Text:   "How many cats did we herd yesterday?",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1525215129, 0),
			OriginalValue: "1525215129.000001",
		},
		ChannelID: "D0PNCRP9N",
		EventTimeStamp: &TimeStamp{
			Time:          time.Unix(1525215129, 0),
			OriginalValue: "1525215129.000001",
		},
		ChannelType: "app_home",
	},
	"message.channels": &MessageChannels{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		ChannelID: "C024BE91L",
		UserID:    "U2147483697",
		Text:      "Live long and prospect.",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1355517523, 0),
			OriginalValue: "1355517523.000005",
		},
		EventTimeStamp: &TimeStamp{
			Time:          time.Unix(1355517523, 0),
			OriginalValue: "1355517523.000005",
		},
		ChannelType: "channel",
	},
	"message.groups": &MessageGroups{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		ChannelID: "G024BE91L",
		UserID:    "U2147483697",
		Text:      "One cannot programmatically detect the difference between `message.mpim` and `message.groups`.",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1355517523, 0),
			OriginalValue: "1355517523.000005",
		},
		EventTimeStamp: &TimeStamp{
			Time:          time.Unix(1355517523, 0),
			OriginalValue: "1355517523.000005",
		},
		ChannelType: "group",
	},
	"message.im": &MessageGroups{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		ChannelID: "D024BE91L",
		UserID:    "U2147483697",
		Text:      "Hello hello can you hear me?",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1355517523, 0),
			OriginalValue: "1355517523.000005",
		},
		EventTimeStamp: &TimeStamp{
			Time:          time.Unix(1355517523, 0),
			OriginalValue: "1355517523.000005",
		},
		ChannelType: "im",
	},
	"message.mpim": &MessageGroups{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		ChannelID: "G024BE91L",
		UserID:    "U2147483697",
		Text:      "Let's make a pact.",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1355517523, 0),
			OriginalValue: "1355517523.000005",
		},
		EventTimeStamp: &TimeStamp{
			Time:          time.Unix(1355517523, 0),
			OriginalValue: "1355517523.000005",
		},
		ChannelType: "mpim",
	},
	"message.bot_message": &MessageBotMessage{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		SubType: "bot_message",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1358877455, 0),
			OriginalValue: "1358877455.000010",
		},
		Text:     "Pushing is the answer",
		BotID:    "BB12033",
		UserName: "github",
		Icon: &BotIcon{
			Image48: "https://slack.com/path/to/hugbot_48.png",
		},
	},
	"message.channel_archive": &MessageChannelArchive{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		SubType: "channel_archive",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1361482916, 0),
			OriginalValue: "1361482916.000003",
		},
		Text: "<U1234|@cal> archived the channel",
		User: "U1234",
	},
	"message.channel_join": &MessageChannelJoin{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		SubType: "channel_join",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1358877458, 0),
			OriginalValue: "1358877458.000011",
		},
		Text:    "<@U2147483828|cal> has joined the channel",
		User:    "U2147483828",
		Inviter: "U123456789",
	},
	"message.channel_leave": &MessageChannelLeave{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		SubType: "channel_leave",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1358877455, 0),
			OriginalValue: "1358877455.000010",
		},
		Text: "<@U2147483828|cal> has left the channel",
		User: "U2147483828",
	},
	"message.channel_name": &MessageChannelName{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		SubType: "channel_name",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1358877455, 0),
			OriginalValue: "1358877455.000010",
		},
		Text:    "<@U2147483828|cal> has renamed the channek from \"random\" to \"watercooler\"",
		User:    "U2147483828",
		OldName: "random",
		Name:    "watercooler",
	},
	"message.channel_topic": &MessageChannelTopic{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		SubType: "channel_topic",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1358877455, 0),
			OriginalValue: "1358877455.000010",
		},
		Text:  "<@U2147483828|cal> set the channel topic: hello world",
		User:  "U2147483828",
		Topic: "hello world",
	},
	"message.channel_unarchive": &MessageChannelUnarchive{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		SubType: "channel_unarchive",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1361482916, 0),
			OriginalValue: "1361482916.000003",
		},
		Text: "<U1234|@cal> un-archived the channel",
		User: "U1234",
	},
	"message.ekm_access_denied": &MessageEKMAccessDenied{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		SubType: "ekm_access_denied",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1358877455, 0),
			OriginalValue: "1358877455.000010",
		},
		Text: "Your admins have suspended everyone's access to this content.",
		User: "UREVOKEDU",
	},
	"message.file_comment": &MessageFileComment{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		SubType: "file_comment",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1361482916, 0),
			OriginalValue: "1361482916.000003",
		},
		Text: "<@cal> commented on a file: ...",
		File: &File{},
		Comment: &Comment{
			ID: "C12345",
			TimeStamp: &TimeStamp{
				Time:          time.Unix(1360782804, 0),
				OriginalValue: "1360782804",
			},
			UserID:  "U1234",
			Content: "comment content",
		},
	},
	"message.file_mention": &MessageFileMention{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		SubType: "file_mention",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1358877455, 0),
			OriginalValue: "1358877455.000010",
		},
		Text: "<@cal> mentioned a file: <https:...7.png|7.png>",
		File: &File{},
		User: "U2147483697",
	},
	"message.group_archive": &MessageGroupArchive{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		SubType: "group_archive",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1361482916, 0),
			OriginalValue: "1361482916.000003",
		},
		Text:    "<U1234|@cal> archived the group",
		User:    "U1234",
		Members: []UserID{"U1234", "U5678"},
	},
	"message.group_join": &MessageGroupJoin{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		SubType: "group_join",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1358877458, 0),
			OriginalValue: "1358877458.000011",
		},
		Text: "<@U2147483828|cal> has joined the group",
		User: "U2147483828",
	},
	"message.group_leave": &MessageGroupLeave{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		SubType: "group_leave",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1358877455, 0),
			OriginalValue: "1358877455.000010",
		},
		Text: "<@U2147483828|cal> has left the group",
		User: "U2147483828",
	},
	"message.group_name": &MessageGroupName{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		SubType: "group_name",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1358877455, 0),
			OriginalValue: "1358877455.000010",
		},
		Text:    "<@U2147483828|cal> has renamed the group from \"random\" to \"watercooler\"",
		User:    "U2147483828",
		OldName: "random",
		Name:    "watercooler",
	},
	"message.group_purpose": &MessageGroupPurpose{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		SubType: "group_purpose",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1358877455, 0),
			OriginalValue: "1358877455.000010",
		},
		Text:    "<@U2147483828|cal> set the group purpose: whatever",
		User:    "U2147483828",
		Purpose: "whatever",
	},
	"message.group_topic": &MessageGroupTopic{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		SubType: "group_topic",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1358877455, 0),
			OriginalValue: "1358877455.000010",
		},
		Text:  "<@U2147483828|cal> set the group topic: hello world",
		User:  "U2147483828",
		Topic: "hello world",
	},
	"message.group_unarchive": &MessageGroupUnarchive{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		SubType: "group_unarchive",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1361482916, 0),
			OriginalValue: "1361482916.000003",
		},
		Text: "<U1234|@cal> un-archived the group",
		User: "U1234",
	},
	"message.me_message": &MessageMeMessage{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		SubType: "me_message",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1355517523, 0),
			OriginalValue: "1355517523.000005",
		},
		Text:    "is doing that thing",
		User:    "U2147483697",
		Channel: "C2147483705",
	},
	"message.message_changed": &MessageMessageChanged{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		SubType: "message_changed",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1358878755, 0),
			OriginalValue: "1358878755.000001",
		},
		Channel: "C2147483705",
		Message: &Message{
			TypedEvent: TypedEvent{
				Type: "message",
			},
			SenderID: "U2147483697",
			Text:     "Hello, world!",
			TimeStamp: &TimeStamp{
				Time:          time.Unix(1355517523, 0),
				OriginalValue: "1355517523.000005",
			},
			Edited: &struct {
				User      UserID     `json:"user"`
				TimeStamp *TimeStamp `json:"ts"`
			}{
				User: "U2147483697",
				TimeStamp: &TimeStamp{
					Time:          time.Unix(1358878755, 0),
					OriginalValue: "1358878755.000001",
				},
			},
		},
		Hidden: true,
	},
	"message.message_deleted": &MessageMessageDeleted{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		SubType: "message_deleted",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1358878755, 0),
			OriginalValue: "1358878755.000001",
		},
		Channel: "C2147483705",
		Hidden:  true,
		DeletedTimeStamp: &TimeStamp{
			Time:          time.Unix(1358878749, 0),
			OriginalValue: "1358878749.000002",
		},
	},
	"message.message_replied": &MessageMessageReplied{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		SubType: "message_replied",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1483037604, 0),
			OriginalValue: "1483037604.017506",
		},
		Channel: "C061EG9SL",
		Hidden:  true,
		EventTimeStamp: &TimeStamp{
			Time:          time.Unix(1483037604, 0),
			OriginalValue: "1483037604.017506",
		},
		Message: &Message{
			TypedEvent: TypedEvent{
				Type: "message",
			},
			SenderID: "U061F7TRS",
			Text:     "Was there was there was there what was there was there what was there was there there was there.",
			TimeStamp: &TimeStamp{
				Time:          time.Unix(1482960137, 0),
				OriginalValue: "1482960137.003543",
			},
			ThreadTimeStamp: &TimeStamp{
				Time:          time.Unix(1482960137, 0),
				OriginalValue: "1482960137.003543",
			},
			Replies: []*struct {
				User      UserID     `json:"user"`
				TimeStamp *TimeStamp `json:"ts"`
			}{
				{
					User: "U061F7AUR",
					TimeStamp: &TimeStamp{
						Time:          time.Unix(1483037603, 0),
						OriginalValue: "1483037603.017503",
					},
				},
			},
		},
	},
	"message.pinned_item": &MessagePinnedItem{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		SubType: "pinned_item",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1360782804, 0),
			OriginalValue: "1360782804.083113",
		},
		User:     "U024BE7LH",
		Text:     "<@U024BE7LH|cal> pinned their Image <https:...7.png|7.png> to this channel.",
		Channel:  "C02ELGNBH",
		ItemType: ItemTypeFile,
		Item:     &Item{},
	},
	"message.thread_broadcast": &MessageThreadBroadcast{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		SubType: "thread_broadcast",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1517414906, 0),
			OriginalValue: "1517414906.000889",
		},
		User: "U061F7AUR",
		Text: "Is but a message within a thread",
		Root: &Message{
			TypedEvent: TypedEvent{
				Type: "message",
			},
			SenderID: "U061F7AUR",
			Text:     "All that we see or seem",
			TimeStamp: &TimeStamp{
				Time:          time.Unix(1517414896, 0),
				OriginalValue: "1517414896.001003",
			},
			ThreadTimeStamp: &TimeStamp{
				Time:          time.Unix(1517414896, 0),
				OriginalValue: "1517414896.001003",
			},
			Replies: []*struct {
				User      UserID     `json:"user"`
				TimeStamp *TimeStamp `json:"ts"`
			}{
				{
					User: "U061F7AUR",
					TimeStamp: &TimeStamp{
						Time:          time.Unix(1517414906, 0),
						OriginalValue: "1517414906.000889",
					},
				},
			},
		},
		ThreadTimeStamp: &TimeStamp{
			Time:          time.Unix(1517414896, 0),
			OriginalValue: "1517414896.001003",
		},
	},
	"message.unpinned_item": &MessageUnpinnedItem{
		TypedEvent: TypedEvent{
			Type: "message",
		},
		SubType: "unpinned_item",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1360782804, 0),
			OriginalValue: "1360782804.083113",
		},
		User:     "USLACKBOT",
		Text:     "<@U024BE7LH|cal> unpinned the message you pinned to the secretplans group.",
		ItemType: ItemTypeGroupMessage,
		Item:     &Item{},
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
	"resources_added": &ResourcesAdded{
		TypedEvent: TypedEvent{
			Type: "resources_added",
		},
		Resources: []*struct {
			Resource *Resource `json:"resource"`
			Scopes   []string  `json:"scopes"`
		}{
			{
				Resource: &Resource{
					Type: "im",
					Grant: &struct {
						Type       string     `json:"type"`
						ResourceID ResourceID `json:"resource_id"`
					}{
						Type:       "specific",
						ResourceID: "DXXXXXXXX",
					},
				},
				Scopes: []string{
					"chat:write:user",
					"im:read",
					"im:history",
					"commands",
				},
			},
		},
	},
	"resources_removed": &ResourcesAdded{
		TypedEvent: TypedEvent{
			Type: "resources_removed",
		},
		Resources: []*struct {
			Resource *Resource `json:"resource"`
			Scopes   []string  `json:"scopes"`
		}{
			{
				Resource: &Resource{
					Type: "im",
					Grant: &struct {
						Type       string     `json:"type"`
						ResourceID ResourceID `json:"resource_id"`
					}{
						Type:       "specific",
						ResourceID: "DXXXXXXXX",
					},
				},
				Scopes: []string{
					"chat:write:user",
					"im:read",
					"im:history",
					"commands",
				},
			},
		},
	},
	"scope_denied": &ScopeDenied{
		TypedEvent: TypedEvent{
			Type: "scope_denied",
		},
		Scopes: []string{
			"files:read",
			"files:write",
			"chat:write",
		},
		TriggerID: "241582872337.47445629121.string",
	},
	"scope_granted": &ScopeDenied{
		TypedEvent: TypedEvent{
			Type: "scope_granted",
		},
		Scopes: []string{
			"files:read",
			"files:write",
			"chat:write",
		},
		TriggerID: "241582872337.47445629121.string",
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
	"tokens_revoked": &TokensRevoked{
		TypedEvent: TypedEvent{
			Type: "tokens_revoked",
		},
		Tokens: &struct {
			OAuth []string `json:"oauth"`
			Bot   []string `json:"bot"`
		}{
			OAuth: []string{"UXXXXXXXX"},
			Bot:   []string{"UXXXXXXXX"},
		},
	},
	"user_change": &UserChanged{
		TypedEvent: TypedEvent{
			Type: "user_change",
		},
		User: &User{},
	},
	"user_resource_denied": &UserResourceDenied{
		TypedEvent: TypedEvent{
			Type: "user_resource_denied",
		},
		UserID: "WXXXXXXXX",
		Scopes: []string{
			"reminders:write:user",
			"reminders:read:user",
		},
		TriggerID: "27082968880.6048553856.5eb9c671f75c636135fdb6bb9e87b606",
	},
	"user_resource_granted": &UserResourceDenied{
		TypedEvent: TypedEvent{
			Type: "user_resource_granted",
		},
		UserID: "WXXXXXXXX",
		Scopes: []string{
			"reminders:write:user",
			"reminders:read:user",
		},
		TriggerID: "27082968880.6048553856.5eb9c671f75c636135fdb6bb9e87b606",
	},
	"user_resource_removed": &UserResourceRemoved{
		TypedEvent: TypedEvent{
			Type: "user_resource_removed",
		},
		UserID:    "WXXXXXXXX",
		TriggerID: "27082968880.6048553856.5eb9c671f75c636135fdb6bb9e87b606",
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

			expected, ok := expectedPayloads[eventName]
			if !ok {
				t.Fatalf("Expected payload for %s is not defined in test", eventName)
			}

			testutil.Compare([]string{eventName}, reflect.ValueOf(expected), reflect.ValueOf(typed), t)
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
		_, err := Decode([]byte(`{
				"type": "UNKNOWN_VALUE"`))
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
