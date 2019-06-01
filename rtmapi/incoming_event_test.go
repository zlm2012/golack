package rtmapi

import (
	"fmt"
	"github.com/oklahomer/golack/slackobject"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

// https://api.slack.com/rtm
// TypedEvent field does not have to be set in each declaration.
// In testStructure(), there is a cheat to set corresponding TypedEvent field.
var expectedPayloads = map[EventType]EventTyper{
	AccountsChangedEvent: &AccountsChanged{},
	BotAddedEvent: &BotAdded{
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
	BotChangedEvent: &BotChanged{
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
	ChannelArchivedEvent: &ChannelArchived{
		ChannelID: "C024BE91L",
		UserID:    "U024BE7LH",
	},
	ChannelCreatedEvent: &ChannelCreated{
		Channel: &struct {
			ID        slackobject.ChannelID `json:"id"`
			Name      string                `json:"name"`
			Created   *TimeStamp            `json:"created"`
			CreatorID slackobject.UserID    `json:"creator"`
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
	ChannelDeletedEvent: &ChannelDeleted{
		ChannelID: "C024BE91L",
	},
	ChannelHistoryChangedEvent: &ChannelHistoryChanged{
		historyChangedEvent: historyChangedEvent{
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
	ChannelJoinedEvent: &ChannelJoined{
		ChannelID: "C024BE91L",
	},
	ChannelLeftEvent: &ChannelLeft{
		ChannelID: "C024BE91L",
	},
	ChannelMarkedEvent: &ChannelMarked{
		markedAsReadEvent: markedAsReadEvent{
			ChannelID: "C024BE91L",
			TimeStamp: &TimeStamp{
				Time:          time.Unix(1401383885, 0),
				OriginalValue: "1401383885.000061",
			},
		},
	},
	ChannelRenameEvent: &ChannelRenamed{
		Channel: &struct {
			ID      slackobject.ChannelID `json:"id"`
			Name    string                `json:"name"`
			Created *TimeStamp            `json:"created"`
		}{
			ID:   "C02ELGNBH",
			Name: "new_name",
			Created: &TimeStamp{
				Time:          time.Unix(1360782804, 0),
				OriginalValue: "1360782804",
			},
		},
	},
	ChannelUnarchiveEvent: &ChannelUnarchived{
		ChannelID: "C024BE91L",
		UserID:    "U024BE7LH",
	},
	CommandsChangedEvent: &CommandsChanged{
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1361482916, 0),
			OriginalValue: "1361482916.000004",
		},
	},
	DNDUpdatedEvent: &DNDUpdated{
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
	DNDUpdatedUserEvent: &DNDUpdatedUser{
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
	EmailDomainChangedEvent: &EmailDomainChanged{
		EmailDomain: "example.com",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1360782804, 0),
			OriginalValue: "1360782804.083113",
		},
	},
	EmojiChangedEvent: &EmojiChanged{
		Subtype: "remove",
		Names:   []string{"picard_facepalm"},
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1361482916, 0),
			OriginalValue: "1361482916.000004",
		},
	},
	FileChangeEvent: &FileChanged{
		FileID: "F2147483862",
		File: &struct {
			ID slackobject.FileID `json:"id"`
		}{
			ID: "F2147483862",
		},
	},
	FileCommentAddedEvent: &FileCommentAdded{
		FileID: "F2147483862",
		File: &struct {
			ID slackobject.FileID `json:"id"`
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
	FileCommentDeletedEvent: &FileCommentDeleted{
		CommentID: "Fc67890",
		FileID:    "F2147483862",
		File: &struct {
			ID slackobject.FileID `json:"id"`
		}{
			ID: "F2147483862",
		},
	},
	FileCommentEditedEvent: &FileCommentEdited{
		FileID: "F2147483862",
		File: &struct {
			ID slackobject.FileID `json:"id"`
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
	FileCreatedEvent: &FileCreated{
		FileID: "F2147483862",
		File: &struct {
			ID slackobject.FileID `json:"id"`
		}{
			ID: "F2147483862",
		},
	},
	FileDeletedEvent: &FileDeleted{
		FileID: "F2147483862",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1361482916, 0),
			OriginalValue: "1361482916.000004",
		},
	},
	FilePublicEvent: &FilePublished{
		FileID: "F2147483862",
		File: &struct {
			ID slackobject.FileID `json:"id"`
		}{
			ID: "F2147483862",
		},
	},
	FileSharedEvent: &FileShared{
		FileID: "F2147483862",
		File: &struct {
			ID slackobject.FileID `json:"id"`
		}{
			ID: "F2147483862",
		},
	},
	FileUnsharedEvent: &FileUnshared{
		FileID: "F2147483862",
		File: &struct {
			ID slackobject.FileID `json:"id"`
		}{
			ID: "F2147483862",
		},
	},
	GoodByeEvent: &GoodBye{},
	GroupArchiveEvent: &GroupArchived{
		ChannelID: "G024BE91L",
	},
	GroupCloseEvent: &GroupClosed{
		UserID:    "U024BE7LH",
		ChannelID: "G024BE91L",
	},
	GroupDeletedEvent: &GroupDeleted{
		ChannelID: "G0QN9RGTT",
	},
	GroupHistoryChangedEvent: &GroupHistoryChanged{
		historyChangedEvent: historyChangedEvent{
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
	GroupJoinedEvent: &GroupJoined{
		ChannelID: "G0QN9RGTT",
	},
	GroupLeftEvent: &GroupLeft{
		ChannelID: "G02ELGNBH",
	},
	GroupMarkedEvent: &GroupMarked{
		markedAsReadEvent: markedAsReadEvent{
			ChannelID: "G024BE91L",
			TimeStamp: &TimeStamp{
				Time:          time.Unix(1401383885, 0),
				OriginalValue: "1401383885.000061",
			},
		},
	},
	GroupOpenEvent: &GroupOpened{
		UserID:    "U024BE7LH",
		ChannelID: "G024BE91L",
	},
	GroupRenameEvent: &GroupRenamed{
		Channel: &struct {
			ID      slackobject.ChannelID `json:"id"`
			Name    string                `json:"name"`
			Created *TimeStamp            `json:"created"`
		}{
			ID:   "G02ELGNBH",
			Name: "new_name",
			Created: &TimeStamp{
				Time:          time.Unix(1360782804, 0),
				OriginalValue: "1360782804",
			},
		},
	},
	GroupUnarchiveEvent: &GroupUnarchived{
		ChannelID: "G024BE91L",
	},
	HelloEvent: &Hello{},
	IMCloseEvent: &IMClosed{
		UserID:    "U024BE7LH",
		ChannelID: "D024BE91L",
	},
	IMCreatedEvent: &IMCreated{
		UserID: "U024BE7LH",
		Channel: &struct {
			ID slackobject.ChannelID `json:"id"`
		}{
			ID: "D024BE91L",
		},
	},
	IMHistoryChangedEvent: &IMHistoryChanged{
		historyChangedEvent: historyChangedEvent{
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
	IMOpenEvent: &IMOpened{
		UserID:    "U024BE7LH",
		ChannelID: "D024BE91L",
	},
	ManualPresenceChangeEvent: &PresenceManuallyChanged{
		Presence: "away",
	},
	MemberJoinedChannelEvent: &MemberJoinedChannel{
		UserID:      "W06GH7XHN",
		ChannelID:   "C0698JE0H",
		ChannelType: "C",
		TeamID:      "T024BE7LD",
		InviterID:   "U123456789",
	},
	MessageEvent: &Message{
		ChannelID: "C2147483705",
		SenderID:  "U2147483697",
		Text:      "Hello world",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1355517523, 0),
			OriginalValue: "1355517523.000005",
		},
	},
	PinAddedEvent: &PinAdded{
		UserID:    "U024BE7LH",
		ChannelID: "C02ELGNBH",
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1360782804, 0),
			OriginalValue: "1360782804.083113",
		},
	},
	PinRemovedEvent: &PinRemoved{
		UserID:    "U024BE7LH",
		ChannelID: "C02ELGNBH",
		HasPins:   false,
		TimeStamp: &TimeStamp{
			Time:          time.Unix(1360782804, 0),
			OriginalValue: "1360782804.083113",
		},
	},
	PrefChangeEvent: &PreferenceChanged{
		Name:  "messages_theme",
		Value: "dense",
	},
	PresenceChangeEvent: &PresenceChange{
		UserID:   "U024BE7LH",
		Presence: "away",
	},
	PresenceQueryEvent: &PresenceQuery{
		UserIDs: []slackobject.UserID{
			"U061F7AUR",
			"W123456",
		},
	},
	PresenceSubEvent: &PresenceSubscribe{
		UserIDs: []slackobject.UserID{
			"U061F7AUR",
			"W123456",
		},
	},
	ReactionAddedEvent: &ReactionAdded{
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
	ReactionRemovedEvent: &ReactionRemoved{
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
	ReconnectURLEvent: &ReconnectURL{},
	StarAddedEvent: &StarAdded{
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
	StarRemovedEvent: &StarRemoved{
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
	SubTeamCreatedEvent: &SubTeamCreated{
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
	SubTeamMembersChangedEvent: &SubTeamMembersChanged{
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
		AddedUserIDs: []slackobject.UserID{
			"U060RNRCZ",
			"U060ULRC0",
			"U061309JM",
		},
		AddedUserCount: 3,
		RemovedUserIDs: []slackobject.UserID{
			"U06129G2V",
		},
		RemovedUsersCount: 1,
	},
	SubTeamSelfAddedEvent: &SubTeamSelfAdded{
		SubTeamID: "S0615G0KT",
	},
	SubTeamSelfRemovedEvent: &SubTeamSelfRemoved{
		SubTeamID: "S0615G0KT",
	},
	SubTeamUpdatedEvent: &SubTeamUpdated{
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
			UserIDs: []slackobject.UserID{
				"U060RNRCZ",
				"U060ULRC0",
				"U06129G2V",
				"U061309JM",
			},
		},
	},
	TeamDomainChangeEvent: &TeamDomainChanged{
		URL:    "https://my.slack.com",
		Domain: "my",
	},
	TeamJoinEvent: &TeamJoined{
		User: &User{},
	},
	TeamMigrationStartedEvent: &TeamMigrationStarted{},
	TeamPlanChangedEvent: &TeamPlanChanged{
		Plan:      "std",
		CanAddUra: false,
		PaidFeatures: []string{
			"feature1",
			"feature2",
		},
	},
	TeamPrefChangeEvent: &TeamPreferenceChanged{
		Name:  "slackbot_responses_only_admins",
		Value: true,
	},
	TeamProfileChangeEvent: &TeamProfileChanged{
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
	TeamProfileDeleteEvent: &TeamProfileDeleted{
		Profile: &struct {
			Fields []string `json:"fields"`
		}{
			Fields: []string{
				"Xf06054AAA",
			},
		},
	},
	TeamProfileReorderEvent: &TeamProfileReordered{
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
	TeamRenameEvent: &TeamRenamed{
		Name: "New Team Name Inc.",
	},
	UserChangeEvent: &UserChanged{
		User: &User{},
	},
	UserTypingEvent: &UserTyping{
		ChannelID: "C02ELGNBH",
		UserID:    "U024BE7LH",
	},
}

func Test_decodePayload_all(t *testing.T) {
	directory := filepath.Join("..", "testdata", "incoming_event")

	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Fatalf("Failed to read testdata directory: %s", err.Error())
		}

		if !strings.HasSuffix(path, ".json") {
			// Skip irrelevant file
			return nil
		}

		filename := filepath.Base(path)
		t.Run(filename, func(t *testing.T) {

			buf, err := ioutil.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to read file: %s. Error: %s.", path, err.Error())
			}

			decodedPayload, err := decodePayload(buf)
			if err != nil {
				t.Fatalf("Failed to decode JSON: %s. Error: %s.", path, err.Error())
			}

			testStructure(decodedPayload, t)
		})

		return nil
	})
}

func testStructure(payload DecodedPayload, t *testing.T) {
	typed, ok := payload.(EventTyper)
	if !ok {
		t.Errorf("Decoded paylaod is not typed event: %#v", payload)
		return
	}

	expected, ok := expectedPayloads[typed.EventType()]
	if !ok {
		t.Fatalf("Expected payload for %s is not defined in test", typed.EventType())
	}

	// A cheat to set TypedEvent field
	reflect.ValueOf(expected).Elem().FieldByName("Type").Set(reflect.ValueOf(typed.EventType()))
	compare([]string{typed.EventType().String()}, reflect.ValueOf(expected), reflect.ValueOf(typed), t)
}

// compare is similar to reflect.DeepEqual, but this prints out more useful information with error message as below:
// [team_profile_reorder > Profile > Fields > Element index at 0 > ID] Expected string is not set. Expected: Xf06054AAA. Actual: INVALID_ID.
// [file_comment_edited > Comment > Content] Expected string is not set. Expected: comment content. Actual: different comment text.
func compare(hierarchy []string, expected reflect.Value, actual reflect.Value, t *testing.T) {
	if expected.Kind() != actual.Kind() {
		t.Errorf("%s Expected type is %s, but is %s.", hierarchyStr(hierarchy), expected.String(), actual.String())
		return
	}

	// Check zero value
	if !expected.IsValid() {
		if actual.IsValid() {
			fmt.Printf("%s is expected to be zero value, but is not: %#v", hierarchyStr(hierarchy), actual.String())
		}
		return
	}

	switch expected.Kind() {
	case reflect.String:
		if expected.String() != actual.String() {
			t.Errorf("%s Expected %s is not set. Expected: %s. Actual: %s.",
				hierarchyStr(hierarchy), expected.Kind().String(), expected.String(), actual.String())
		}

	case reflect.Bool:
		if expected.Bool() != actual.Bool() {
			t.Errorf("%s Expected %s is not set. Expected: %t. Actual: %t",
				hierarchyStr(hierarchy), expected.Kind().String(), expected.Bool(), actual.Bool())
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if expected.Uint() != actual.Uint() {
			t.Errorf(
				"%s Expected %s is not set. Expected: %d. Actual: %d.",
				hierarchyStr(hierarchy), expected.Kind().String(), expected.Uint(), actual.Uint())
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if expected.Int() != actual.Int() {
			t.Errorf(
				"%s Expected %s is not set. Expected: %d. Actual: %d.",
				hierarchyStr(hierarchy), expected.Kind().String(), expected.Int(), actual.Int())
		}

	case reflect.Ptr:
		compare(hierarchy, expected.Elem(), actual.Elem(), t)

	case reflect.Struct:
		for i := 0; i < expected.NumField(); i++ {
			tmp := make([]string, len(hierarchy))
			copy(tmp, hierarchy)
			tmp = append(tmp, expected.Type().Field(i).Name)
			compare(tmp, expected.Field(i), actual.Field(i), t)
		}

	case reflect.Array, reflect.Slice:
		if expected.Len() != actual.Len() {
			t.Errorf("%s Element %s size differs. Expected: %d. Actual: %d.",
				hierarchyStr(hierarchy), expected.Kind().String(), expected.Len(), actual.Len())
			return
		}

		tmp := make([]string, len(hierarchy))
		copy(tmp, hierarchy)
		for i := 0; i < expected.Len(); i++ {
			tmp = append(tmp, fmt.Sprintf("Element index at %d", i))
			compare(tmp, expected.Index(i), actual.Index(i), t)
		}

	default:
		t.Errorf("Uncontrollable Kind %s is given: %#v", expected.Kind().String(), expected)

	}
}

func hierarchyStr(stack []string) string {
	return fmt.Sprintf("[%s]", strings.Join(stack, " > "))
}
