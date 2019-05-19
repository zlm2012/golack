package rtmapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/oklahomer/golack/slackobject"
	"github.com/tidwall/gjson"
)

var (
	ErrEmptyPayload = errors.New("empty payload was given")
)

type UnexpectedMessageTypeError struct {
	MessageType int
	Payload     []byte
}

func (e *UnexpectedMessageTypeError) Error() string {
	return fmt.Sprintf("unexpected message type, %d, is given: %s", e.MessageType, e.Payload)
}

type DecodedPayload interface{}

type PayloadReceiver interface {
	Receive() (DecodedPayload, error)
}

type PayloadSender interface {
	Send(slackobject.ChannelID, string) error
	Ping() error
}

type Connection interface {
	PayloadReceiver
	PayloadSender
	io.Closer
}

// Connect connects to Slack WebSocket server.
func Connect(_ context.Context, url string) (Connection, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
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
	messageType, payload, err := wrapper.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	// Only TextMessage is supported by RTM API.
	if messageType != websocket.TextMessage {
		return nil, &UnexpectedMessageTypeError{MessageType: messageType, Payload: payload}
	}

	decoded, err := decodePayload(payload)
	return decoded, err
}

func (wrapper *connWrapper) Send(channel slackobject.ChannelID, content string) error {
	event := NewOutgoingMessage(wrapper.outgoingEventID, channel, content)
	return wrapper.conn.WriteJSON(event)
}

func (wrapper *connWrapper) Ping() error {
	ping := NewPing(wrapper.outgoingEventID)
	return wrapper.conn.WriteJSON(ping)
}

func (wrapper *connWrapper) Close() error {
	return wrapper.conn.Close()
}

var (
	eventTypeMap = map[EventType]reflect.Type{
		AccountsChangedEvent:       reflect.TypeOf(&AccountsChanged{}).Elem(),
		BotAddedEvent:              reflect.TypeOf(&BotAdded{}).Elem(),
		BotChangedEvent:            reflect.TypeOf(&BotChanged{}).Elem(),
		ChannelArchivedEvent:       reflect.TypeOf(&ChannelArchived{}).Elem(),
		ChannelCreatedEvent:        reflect.TypeOf(&ChannelCreated{}).Elem(),
		ChannelDeletedEvent:        reflect.TypeOf(&ChannelDeleted{}).Elem(),
		ChannelHistoryChangedEvent: reflect.TypeOf(&ChannelHistoryChanged{}).Elem(),
		ChannelJoinedEvent:         reflect.TypeOf(&ChannelJoined{}).Elem(),
		ChannelLeftEvent:           reflect.TypeOf(&ChannelLeft{}).Elem(),
		ChannelMarkedEvent:         reflect.TypeOf(&ChannelMarked{}).Elem(),
		ChannelRenameEvent:         reflect.TypeOf(&ChannelRenamed{}).Elem(),
		ChannelUnarchiveEvent:      reflect.TypeOf(&ChannelUnarchived{}).Elem(),
		CommandsChangedEvent:       reflect.TypeOf(&CommandsChanged{}).Elem(),
		DNDUpdatedEvent:            reflect.TypeOf(&DNDUpdated{}).Elem(),
		DNDUpdatedUserEvent:        reflect.TypeOf(&DNDUpdated{}).Elem(),
		EmailDomainChangedEvent:    reflect.TypeOf(&EmailDomainChanged{}).Elem(),
		EmojiChangedEvent:          reflect.TypeOf(&EmojiChanged{}).Elem(),
		FileChangeEvent:            reflect.TypeOf(&FileChanged{}).Elem(),
		FileCommentAddedEvent:      reflect.TypeOf(&FileCommentAdded{}).Elem(),
		FileCommentDeletedEvent:    reflect.TypeOf(&FileCommentDeleted{}).Elem(),
		FileCommentEditedEvent:     reflect.TypeOf(&FileCommentEdited{}).Elem(),
		FileCreatedEvent:           reflect.TypeOf(&FileCreated{}).Elem(),
		FileDeletedEvent:           reflect.TypeOf(&FileDeleted{}).Elem(),
		FilePublicEvent:            reflect.TypeOf(&FilePublicated{}).Elem(),
		FileSharedEvent:            reflect.TypeOf(&FileShared{}).Elem(),
		FileUnsharedEvent:          reflect.TypeOf(&FileUnshared{}).Elem(),
		GoodByeEvent:               reflect.TypeOf(&GoodBye{}).Elem(),
		GroupArchiveEvent:          reflect.TypeOf(&GroupArchived{}).Elem(),
		GroupCloseEvent:            reflect.TypeOf(&GroupClosed{}).Elem(),
		GroupHistoryChangedEvent:   reflect.TypeOf(&GroupHistoryChanged{}).Elem(),
		GroupJoinedEvent:           reflect.TypeOf(&GroupJoined{}).Elem(),
		GroupLeftEvent:             reflect.TypeOf(&GroupLeft{}).Elem(),
		GroupMarkedEvent:           reflect.TypeOf(&GroupMarked{}).Elem(),
		GroupOpenEvent:             reflect.TypeOf(&GroupOpened{}).Elem(),
		GroupRenameEvent:           reflect.TypeOf(&GroupRenamed{}).Elem(),
		GroupUnarchiveEvent:        reflect.TypeOf(&GroupUnarchived{}).Elem(),
		HelloEvent:                 reflect.TypeOf(&Hello{}).Elem(),
		IMCloseEvent:               reflect.TypeOf(&IMClosed{}).Elem(),
		IMCreatedEvent:             reflect.TypeOf(&IMCreated{}).Elem(),
		IMHistoryChangedEvent:      reflect.TypeOf(&IMHistoryChanged{}).Elem(),
		IMMarkedEvent:              reflect.TypeOf(&IMMarked{}).Elem(),
		IMOpenEvent:                reflect.TypeOf(&IMOpened{}).Elem(),
		ManualPresenceChangeEvent:  reflect.TypeOf(&PresenceManuallyChanged{}).Elem(),
		PinAddedEvent:              reflect.TypeOf(&PinAdded{}).Elem(),
		PinRemovedEvent:            reflect.TypeOf(&PinRemoved{}).Elem(),
		PrefChangeEvent:            reflect.TypeOf(&PreferenceChanged{}).Elem(),
		PresenceChangeEvent:        reflect.TypeOf(&PresenceChange{}).Elem(),
		ReactionAddedEvent:         reflect.TypeOf(&ReactionAdded{}).Elem(),
		ReactionRemovedEvent:       reflect.TypeOf(&ReactionRemoved{}).Elem(),
		ReconnectURLEvent:          reflect.TypeOf(&ReconnectURL{}).Elem(),
		StarAddedEvent:             reflect.TypeOf(&StarAdded{}).Elem(),
		StarRemovedEvent:           reflect.TypeOf(&StarRemoved{}).Elem(),
		SubTeamCreatedEvent:        reflect.TypeOf(&SubTeamCreated{}).Elem(),
		SubTeamSlefAddedEvent:      reflect.TypeOf(&SubTeamSelfAdded{}).Elem(),
		SubTeamSelfRemovedEvent:    reflect.TypeOf(&SubTeamSelfRemoved{}).Elem(),
		SubTeamUpdatedEvent:        reflect.TypeOf(&SubTeamUpdated{}).Elem(),
		TeamDomainChangeEvent:      reflect.TypeOf(&TeamDomainChanged{}).Elem(),
		TeamJoinEvent:              reflect.TypeOf(&TeamJoined{}).Elem(),
		TeamMigrationStartedEvent:  reflect.TypeOf(&TeamMigrationStarted{}).Elem(),
		TeamPlanChangedEvent:       reflect.TypeOf(&TeamPlanChanged{}).Elem(),
		TeamPrefChangeEvent:        reflect.TypeOf(&TeamPreferenceChanged{}).Elem(),
		TeamProfileChangeEvent:     reflect.TypeOf(&TeamProfileChanged{}).Elem(),
		TeamProfileDeleteEvent:     reflect.TypeOf(&TeamProfileDeleted{}).Elem(),
		TeamProfileReorderEvent:    reflect.TypeOf(&TeamProfileReordered{}).Elem(),
		TeamRenameEvent:            reflect.TypeOf(&TeamRenamed{}).Elem(),
		UserChangeEvent:            reflect.TypeOf(&UserChanged{}).Elem(),
		UserTypingEvent:            reflect.TypeOf(&UserTyping{}).Elem(),
		PongEvent:                  reflect.TypeOf(&Pong{}).Elem(),
	}
	subTypeMap = map[SubType]reflect.Type{
		BotMessage:       reflect.TypeOf(&MiscMessage{}).Elem(),
		ChannelArchive:   reflect.TypeOf(&MiscMessage{}).Elem(),
		ChannelJoin:      reflect.TypeOf(&MiscMessage{}).Elem(),
		ChannelLeave:     reflect.TypeOf(&MiscMessage{}).Elem(),
		ChannelName:      reflect.TypeOf(&MiscMessage{}).Elem(),
		ChannelPurpose:   reflect.TypeOf(&MiscMessage{}).Elem(),
		ChannelTopic:     reflect.TypeOf(&MiscMessage{}).Elem(),
		ChannelUnarchive: reflect.TypeOf(&MiscMessage{}).Elem(),
		FileComment:      reflect.TypeOf(&MiscMessage{}).Elem(),
		FileMention:      reflect.TypeOf(&MiscMessage{}).Elem(),
		FileShare:        reflect.TypeOf(&MiscMessage{}).Elem(),
		GroupArchive:     reflect.TypeOf(&MiscMessage{}).Elem(),
		GroupJoin:        reflect.TypeOf(&MiscMessage{}).Elem(),
		GroupLeave:       reflect.TypeOf(&MiscMessage{}).Elem(),
		GroupName:        reflect.TypeOf(&MiscMessage{}).Elem(),
		GroupPurpose:     reflect.TypeOf(&MiscMessage{}).Elem(),
		GroupTopic:       reflect.TypeOf(&MiscMessage{}).Elem(),
		GroupUnarchive:   reflect.TypeOf(&MiscMessage{}).Elem(),
		MeMessage:        reflect.TypeOf(&MiscMessage{}).Elem(),
		MessageChanged:   reflect.TypeOf(&MiscMessage{}).Elem(),
		MessageDeleted:   reflect.TypeOf(&MiscMessage{}).Elem(),
		PinnedItem:       reflect.TypeOf(&MiscMessage{}).Elem(),
		UnpinnedItem:     reflect.TypeOf(&MiscMessage{}).Elem(),
	}
	messageEventType          = reflect.TypeOf(&Message{}).Elem()
	websocketReplyOKEventType = reflect.TypeOf(&WebSocketOKReply{}).Elem()
	websocketReplyNGEventType = reflect.TypeOf(&WebSocketNGReply{}).Elem()
)

func decodePayload(input json.RawMessage) (DecodedPayload, error) {
	inputStr := strings.TrimSpace(string(input))
	if len(inputStr) == 0 {
		return nil, ErrEmptyPayload
	}

	res := gjson.Parse(inputStr)
	eventTypeValue := res.Get("type")
	if eventTypeValue.Exists() {
		eventType := AtoEventType(eventTypeValue.String())

		// See if corresponding type can be found from type value.
		mapping, ok := eventTypeMap[eventType]
		if ok {
			return unmarshal(input, mapping)
		}

		// If type value is that of MessageEvent, subtype value must be checked to find corresponding type.
		if eventType == MessageEvent {
			subType := res.Get("subtype")
			if subType.Exists() {
				mapping, ok := subTypeMap[AtoSubType(subType.String())]
				if ok {
					return unmarshal(input, mapping)
				}

				// TODO handle unsupported subtype.
			}

			// If subtype is not given, then this is a plain message.
			return unmarshal(input, messageEventType)
		}

		// event type is given, but there is no matching type.
		return nil, NewMalformedPayloadError(fmt.Sprintf(`unsupported event type "%s" is given: %s`, eventTypeValue.String(), input))
	}

	// https://api.slack.com/rtm#handling_responses
	if res.Get("reply_to").Exists() && res.Get("ok").Exists() {
		if res.Get("ok").Bool() {
			return unmarshal(input, websocketReplyOKEventType)
		}
		return unmarshal(input, websocketReplyNGEventType)
	}

	return nil, NewMalformedPayloadError(fmt.Sprintf("given json object has unknown structure. can not handle: %s.", input))
}

func unmarshal(input json.RawMessage, mapping reflect.Type) (DecodedPayload, error) {
	payload := reflect.New(mapping).Interface()
	err := json.Unmarshal(input, &payload)
	if err != nil {
		return nil, NewMalformedPayloadError(err.Error())
	}

	return payload, nil
}
