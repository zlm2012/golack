package rtmapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
	"io"
	"strings"
)

var (
	ErrEmptyPayload = errors.New("empty payload was given")
)

type DecodedPayload interface{}

type PayloadReceiver interface {
	Receive() (DecodedPayload, error)
}

type PayloadSender interface {
	Send(ChannelID, string) error
	Ping() error
}

type Connection interface {
	PayloadReceiver
	PayloadSender
	io.Closer
}

// Connect connects to Slack WebSocket server.
func Connect(_ context.Context, url string) (Connection, error) {
	conn, err := websocket.Dial(url, "", "http://localhost")
	if err != nil {
		return nil, err
	}

	return newConnectionWrapper(conn), nil
}

// connWrapper is a thin wrapper that wraps WebSocket connection and its methods.
// This instance is created per-connection.
type connWrapper struct {
	conn *websocket.Conn

	// https://api.slack.com/rtm#sending_messages
	// Every event should have a unique (for that connection) positive integer ID.
	outgoingEventID *OutgoingEventID
}

func newConnectionWrapper(conn *websocket.Conn) Connection {
	return &connWrapper{
		conn:            conn,
		outgoingEventID: NewOutgoingEventID(),
	}

}

// Receive is a blocking method to receive payload from WebSocket connection.
// When connection is closed in the middle of this method call, this immediately returns error.
func (wrapper *connWrapper) Receive() (DecodedPayload, error) {
	// Slack's RTM events and reply all have different form
	// so websocket.JSON.Receive can only work with json.RawMessage,
	// which ends up with multiple json.Unmarshal calls for proper mapping.
	payload := []byte{}
	err := websocket.Message.Receive(wrapper.conn, &payload)
	if err != nil {
		return nil, err
	}

	decoded, err := decodePayload(payload)
	return decoded, err
}

func (wrapper *connWrapper) Send(channel ChannelID, content string) error {
	event := NewOutgoingMessage(wrapper.outgoingEventID, channel, content)
	return websocket.JSON.Send(wrapper.conn, event)
}

func (wrapper *connWrapper) Ping() error {
	ping := NewPing(wrapper.outgoingEventID)
	return websocket.JSON.Send(wrapper.conn, ping)
}

func (wrapper *connWrapper) Close() error {
	return wrapper.conn.Close()
}

func decodePayload(input json.RawMessage) (DecodedPayload, error) {
	inputStr := strings.TrimSpace(string(input))
	if len(inputStr) == 0 {
		return nil, ErrEmptyPayload
	}

	res := gjson.Parse(inputStr)
	eventType := res.Get("type")
	if eventType.Exists() {
		typeVal := eventType.String()
		switch typeVal {
		case AccountsChangedEvent:
			return unmarshal(input, &AccountsChanged{})
		case BotAddedEvent:
			return unmarshal(input, &BotAdded{})
		case BotChangedEvent:
			return unmarshal(input, &BotChanged{})
		case ChannelArchivedEvent:
			return unmarshal(input, &ChannelArchived{})
		case ChannelCreatedEvent:
			return unmarshal(input, &ChannelCreated{})
		case ChannelDeletedEvent:
			return unmarshal(input, &ChannelDeleted{})
		case ChannelHistoryChangedEvent:
			return unmarshal(input, &ChannelHistoryChanged{})
		case ChannelJoinedEvent:
			return unmarshal(input, &ChannelJoined{})
		case ChannelLeftEvent:
			return unmarshal(input, &ChannelLeft{})
		case ChannelMarkedEvent:
			return unmarshal(input, &ChannelMarked{})
		case ChannelRenameEvent:
			return unmarshal(input, &ChannelRenamed{})
		case ChannelUnarchiveEvent:
			return unmarshal(input, &ChannelUnarchived{})
		case CommandsChangedEvent:
			return unmarshal(input, &CommandsChanged{})
		case DNDUpdatedEvent:
			return unmarshal(input, &DNDUpdated{})
		case DNDUpdatedUserEvent:
			return unmarshal(input, &DNDUpdated{})
		case EmailDomainChangedEvent:
			return unmarshal(input, &EmailDomainChanged{})
		case EmojiChangedEvent:
			return unmarshal(input, &EmojiChanged{})
		case FileChangeEvent:
			return unmarshal(input, &FileChanged{})
		case FileCommentAddedEvent:
			return unmarshal(input, &FileCommentAdded{})
		case FileCommentDeletedEvent:
			return unmarshal(input, &FileCommentDeleted{})
		case FileCommentEditedEvent:
			return unmarshal(input, &FileCommentEdited{})
		case FileCreatedEvent:
			return unmarshal(input, &FileCreated{})
		case FileDeletedEvent:
			return unmarshal(input, &FileDeleted{})
		case FilePublicEvent:
			return unmarshal(input, &FilePublicated{})
		case FileSharedEvent:
			return unmarshal(input, &FileShared{})
		case FileUnsharedEvent:
			return unmarshal(input, &FileUnshared{})
		case GoodByeEvent:
			return unmarshal(input, &GoodBye{})
		case GroupArchiveEvent:
			return unmarshal(input, &GroupArchived{})
		case GroupCloseEvent:
			return unmarshal(input, &GroupClosed{})
		case GroupHistoryChangedEvent:
			return unmarshal(input, &GroupHistoryChanged{})
		case GroupJoinedEvent:
			return unmarshal(input, &GroupJoined{})
		case GroupLeftEvent:
			return unmarshal(input, &GroupLeft{})
		case GroupMarkedEvent:
			return unmarshal(input, &GroupMarked{})
		case GroupOpenEvent:
			return unmarshal(input, &GroupOpened{})
		case GroupRenameEvent:
			return unmarshal(input, &GroupRenamed{})
		case GroupUnarchiveEvent:
			return unmarshal(input, &GroupUnarchived{})
		case HelloEvent:
			return unmarshal(input, &Hello{})
		case IMCloseEvent:
			return unmarshal(input, &IMClosed{})
		case IMCreatedEvent:
			return unmarshal(input, &IMCreated{})
		case IMHistoryChangedEvent:
			return unmarshal(input, &IMHistoryChanged{})
		case IMMarkedEvent:
			return unmarshal(input, &IMMarked{})
		case IMOpenEvent:
			return unmarshal(input, &IMOpened{})
		case ManualPresenceChangeEvent:
			return unmarshal(input, &PresenceManuallyChanged{})
		case MessageEvent:
			subType := res.Get("subtype")
			if subType.Exists() {
				switch subType.String() {
				// TODO handle each subtypes
				case BotMessage:
					return unmarshal(input, &MiscMessage{})
				case ChannelArchive:
					return unmarshal(input, &MiscMessage{})
				case ChannelJoin:
					return unmarshal(input, &MiscMessage{})
				case ChannelLeave:
					return unmarshal(input, &MiscMessage{})
				case ChannelName:
					return unmarshal(input, &MiscMessage{})
				case ChannelPurpose:
					return unmarshal(input, &MiscMessage{})
				case ChannelTopic:
					return unmarshal(input, &MiscMessage{})
				case ChannelUnarchive:
					return unmarshal(input, &MiscMessage{})
				case FileComment:
					return unmarshal(input, &MiscMessage{})
				case FileMention:
					return unmarshal(input, &MiscMessage{})
				case FileShare:
					return unmarshal(input, &MiscMessage{})
				case GroupArchive:
					return unmarshal(input, &MiscMessage{})
				case GroupJoin:
					return unmarshal(input, &MiscMessage{})
				case GroupLeave:
					return unmarshal(input, &MiscMessage{})
				case GroupName:
					return unmarshal(input, &MiscMessage{})
				case GroupPurpose:
					return unmarshal(input, &MiscMessage{})
				case GroupTopic:
					return unmarshal(input, &MiscMessage{})
				case GroupUnarchive:
					return unmarshal(input, &MiscMessage{})
				case MeMessage:
					return unmarshal(input, &MiscMessage{})
				case MessageChanged:
					return unmarshal(input, &MiscMessage{})
				case MessageDeleted:
					return unmarshal(input, &MiscMessage{})
				case PinnedItem:
					return unmarshal(input, &MiscMessage{})
				case UnpinnedItem:
					return unmarshal(input, &MiscMessage{})
				default:
					return unmarshal(input, &MiscMessage{})
				}
			} else {
				// plain message
				return unmarshal(input, &Message{})
			}
		case PinAddedEvent:
			return unmarshal(input, &PinAdded{})
		case PinRemovedEvent:
			return unmarshal(input, &PinRemoved{})
		case PrefChangeEvent:
			return unmarshal(input, &PreferenceChanged{})
		case PresenceChangeEvent:
			return unmarshal(input, &PresenceChange{})
		case ReactionAddedEvent:
			return unmarshal(input, &ReactionAdded{})
		case ReactionRemovedEvent:
			return unmarshal(input, &ReactionRemoved{})
		case ReconnectURLEvent:
			return unmarshal(input, &ReconnectURL{})
		case StarAddedEvent:
			return unmarshal(input, &StarAdded{})
		case StarRemovedEvent:
			return unmarshal(input, &StarRemoved{})
		case SubTeamCreatedEvent:
			return unmarshal(input, &SubTeamCreated{})
		case SubTeamSlefAddedEvent:
			return unmarshal(input, &SubTeamSelfAdded{})
		case SubTeamSelfRemovedEvent:
			return unmarshal(input, &SubTeamSelfRemoved{})
		case SubTeamUpdatedEvent:
			return unmarshal(input, &SubTeamUpdated{})
		case TeamDomainChangeEvent:
			return unmarshal(input, &TeamDomainChanged{})
		case TeamJoinEvent:
			return unmarshal(input, &TeamJoined{})
		case TeamMigrationStartedEvent:
			return unmarshal(input, &TeamMigrationStarted{})
		case TeamPlanChangedEvent:
			return unmarshal(input, &TeamPlanChanged{})
		case TeamPrefChangeEvent:
			return unmarshal(input, &TeamPreferenceChanged{})
		case TeamProfileChangeEvent:
			return unmarshal(input, &TeamProfileChanged{})
		case TeamProfileDeleteEvent:
			return unmarshal(input, &TeamProfileDeleted{})
		case TeamProfileReorderEvent:
			return unmarshal(input, &TeamProfileReordered{})
		case TeamRenameEvent:
			return unmarshal(input, &TeamRenamed{})
		case UserChangeEvent:
			return unmarshal(input, &UserChanged{})
		case UserTypingEvent:
			return unmarshal(input, &UserTyping{})
		case PongEvent:
			return unmarshal(input, &Pong{})
		default:
			return nil, NewMalformedPayloadError(fmt.Sprintf(`unsupported event type "%s" is given: %s`, typeVal, input))
		}
	}

	if res.Get("reply_to").Exists() && res.Get("ok").Exists() {
		// https://api.slack.com/rtm#handling_responses
		// When incoming object can't be treated as Event, try treat this as WebSocketReply.
		return unmarshal(input, &WebSocketReply{})
	}

	return nil, NewMalformedPayloadError(fmt.Sprintf("given json object has unknown structure. can not handle: %s.", input))
}

func unmarshal(input json.RawMessage, mapping interface{}) (DecodedPayload, error) {
	if err := json.Unmarshal(input, mapping); err != nil {
		return nil, NewMalformedPayloadError(err.Error())
	}

	return mapping, nil
}
