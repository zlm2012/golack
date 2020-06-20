package event

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"golang.org/x/xerrors"
)

// List of block elements: https://api.slack.com/reference/block-kit/block-elements

func UnmarshalBlockElement(input json.RawMessage) (BlockElement, error) {
	parsed := gjson.ParseBytes(input)
	typeValue := parsed.Get("type")

	var typed BlockElement
	switch t := typeValue.String(); t {
	case "button":
		typed = &ButtonBlockElement{}

	case "checkboxes":
		typed = &CheckboxBlockElement{}

	case "datepicker":
		typed = &DatePickerBlockElement{}

	case "image":
		typed = &ImageBlockElement{}

	case "multi_static_select":
		typed = &MultiStaticSelectBLockElement{}

	case "multi_external_select":
		typed = &MultiExternalSelectBLockElement{}

	case "multi_users_select":
		typed = &MultiUsersSelectBlockElement{}

	case "multi_conversations_select":
		typed = &MultiConversationsSelectBlockElement{}

	case "multi_channels_select":
		typed = &MultiChannelsSelectBlockElement{}

	case "overflow":
		typed = &OverflowBlockElement{}

	case "plain_text_input":
		typed = &PlainTextInputBlockElement{}

	case "radio_buttons":
		typed = &RadioButtonGroupBlockElement{}

	case "static_select":
		typed = &StaticSelectBlockElement{}

	case "external_select":
		typed = &ExternalSelectBlockElement{}

	case "users_select":
		typed = &UsersSelectBlockElement{}

	case "conversations_select":
		typed = &ConversationsSelectBlockElement{}

	case "channels_select":
		typed = &ChannelsSelectBlockElement{}

	default:
		return nil, fmt.Errorf("failed to handle unknown block type: %s", t)
	}

	err := json.Unmarshal(input, typed)
	if err != nil {
		return nil, xerrors.Errorf("failed to unmarshal %T", typed)
	}
	return typed, nil
}

type BlockElement interface {
	BlockElementType() string
}

type blockElement struct {
	Type string `json:"type"`
}

var _ BlockElement = (*blockElement)(nil)

func (be *blockElement) BlockElementType() string {
	return be.Type
}

type ButtonBlockElement struct {
	blockElement
	Text     *TextCompositionObject    `json:"text"`
	ActionID ActionID                  `json:"action_id"`
	URL      string                    `json:"url,omitempty"`
	Value    string                    `json:"value,omitempty"`
	Style    string                    `json:"style,omitempty"` // TODO "primary", "danger" or empty
	Confirm  *ConfirmationDialogObject `json:"confirm,omitempty"`
}

type CheckboxBlockElement struct {
	blockElement
	ActionID       ActionID                  `json:"action_id"`
	Options        []*OptionObject           `json:"options"`
	InitialOptions []*OptionObject           `json:"initial_options,omitempty"`
	Confirm        *ConfirmationDialogObject `json:"confirm,omitempty"`
}

type DatePickerBlockElement struct {
	blockElement
	ActionID    ActionID                  `json:"action_id"`
	PlaceHolder *TextCompositionObject    `json:"placeholder,omitempty"`
	InitialDate string                    `json:"initial_date,omitempty"`
	Confirm     *ConfirmationDialogObject `json:"confirm,omitempty"`
}

type ImageBlockElement struct {
	blockElement
	ImageURL string `json:"image_url"`
	AltText  string `json:"alt_text"`
}

// MultiStaticSelectBLockElement is a Multi-select menu with static options.
// This has a type of "multi_static_select."
type MultiStaticSelectBLockElement struct {
	blockElement
	Placeholder      *TextCompositionObject    `json:"placeholder"`
	ActionID         ActionID                  `json:"action_id"`
	Options          []*OptionObject           `json:"options"`
	OptionsGroups    []*OptionGroupObject      `json:"option_groups,omitempty"`
	InitialOptions   []*OptionObject           `json:"initial_options,omitempty"`
	Confirm          *ConfirmationDialogObject `json:"confirm,omitempty"`
	MaxSelectedItems int                       `json:"max_selected_items,omitempty"`
}

// MultiExternalSelectBLockElement is a Multi-select menu with external data source.
// This has a type of "multi_external_select."
type MultiExternalSelectBLockElement struct {
	blockElement
	Placeholder      *TextCompositionObject    `json:"placeholder"`
	ActionID         ActionID                  `json:"action_id"`
	MinQueryLength   int                       `json:"min_query_length,omitempty"`
	InitialOptions   []*OptionObject           `json:"initial_options,omitempty"`
	Confirm          *ConfirmationDialogObject `json:"confirm,omitempty"`
	MaxSelectedItems int                       `json:"max_selected_items,omitempty"`
}

// MultiUsersSelectBlockElement is a Multi-select menu with a list of Slack users visible to the current user.
// This has a type of "multi_users_select."
type MultiUsersSelectBlockElement struct {
	blockElement
	Placeholder      *TextCompositionObject    `json:"placeholder"`
	ActionID         ActionID                  `json:"action_id"`
	InitialUserIDs   []UserID                  `json:"users,omitempty"`
	Confirm          *ConfirmationDialogObject `json:"confirm,omitempty"`
	MaxSelectedItems int                       `json:"max_selected_items,omitempty"`
}

// MultiConversationsSelectBlockElement is a Multi-select menu with a list of public and private channels.
// This has a type of "multi_conversations_select."
type MultiConversationsSelectBlockElement struct {
	blockElement
	Placeholder                  *TextCompositionObject    `json:"placeholder"`
	ActionID                     ActionID                  `json:"action_id"`
	InitialConversations         []string                  `json:"initial_conversations,omitempty"`
	DefaultToCurrentConversation bool                      `json:"default_to_current_conversation,omitempty"`
	Confirm                      *ConfirmationDialogObject `json:"confirm,omitempty"`
	MaxSelectedItems             int                       `json:"max_selected_items,omitempty"`
	Filter                       *FilterObject             `json:"filter,omitempty"`
}

// MultiChannelsSelectBlockElement is a Multi-select menu with a list of public channels visible to the current user.
// This has a type of "multi_channels_select."
type MultiChannelsSelectBlockElement struct {
	blockElement
	Placeholder       *TextCompositionObject    `json:"placeholder"`
	ActionID          ActionID                  `json:"action_id"`
	InitialChannelIDs []ChannelID               `json:"initial_channels,omitempty"`
	Confirm           *ConfirmationDialogObject `json:"confirm,omitempty"`
	MaxSelectedItems  int                       `json:"max_selected_items,omitempty"`
}

type OverflowBlockElement struct {
	blockElement
	ActionID ActionID                  `json:"action_id"`
	Options  []*OptionObject           `json:"options"`
	Confirm  *ConfirmationDialogObject `json:"confirm"`
}

type PlainTextInputBlockElement struct {
	blockElement
	ActionID     ActionID               `json:"action_id"`
	Placeholder  *TextCompositionObject `json:"placeholder,omitempty"`
	InitialValue string                 `json:"initial_value,omitempty"`
	Multiline    bool                   `json:"multiline,omitempty"`
	MinLength    int                    `json:"min_length,omitempty"`
	MaxLength    int                    `json:"max_length,omitempty"`
}

type RadioButtonGroupBlockElement struct {
	blockElement
	ActionID      ActionID                  `json:"action_id"`
	Options       []*OptionObject           `json:"options"`
	InitialOption *OptionObject             `json:"initial_option,omitempty"`
	Confirm       *ConfirmationDialogObject `json:"confirm,omitempty"`
}

// StaticSelectBlockElement is a Select menu element with a static list of options.
// This has a type of "static_select."
type StaticSelectBlockElement struct {
	blockElement
	Placeholder   *TextCompositionObject    `json:"placeholder"`
	ActionID      ActionID                  `json:"action_id"`
	Options       []*OptionObject           `json:"options"`
	OptionGroups  []*OptionGroupObject      `json:"option_groups,omitempty"`
	InitialOption *OptionObject             `json:"initial_option,omitempty"`
	Confirm       *ConfirmationDialogObject `json:"confirm,omitempty"`
}

// ExternalSelectBlockElement is a Select menu element with external data source.
// This has a type of "external_select."
type ExternalSelectBlockElement struct {
	blockElement
	Placeholder    *TextCompositionObject    `json:"placeholder"`
	ActionID       ActionID                  `json:"action_id"`
	InitialOption  *OptionObject             `json:"initial_option,omitempty"`
	MinQueryLength int                       `json:"min_query_length,omitempty"`
	Confirm        *ConfirmationDialogObject `json:"confirm,omitempty"`
}

// UsersSelectBlockElement is a Select menu element with a static list of Slack users visible to the current user.
// This has a type of "users_select."
type UsersSelectBlockElement struct {
	blockElement
	Placeholder   *TextCompositionObject    `json:"placeholder"`
	ActionID      ActionID                  `json:"action_id"`
	InitialUserID UserID                    `json:"users,omitempty"`
	Confirm       *ConfirmationDialogObject `json:"confirm,omitempty"`
}

// ConversationsSelectBlockElement is a Select menu element with a list of public and private channels, DMs and MPIMs visible to the current user.
// This has a type of "conversations_select."
type ConversationsSelectBlockElement struct {
	blockElement
	Placeholder                  *TextCompositionObject    `json:"placeholder"`
	ActionID                     ActionID                  `json:"action_id"`
	InitialConversation          string                    `json:"initial_conversation,omitempty"`
	DefaultToCurrentConversation bool                      `json:"default_to_current_conversation,omitempty"`
	Confirm                      *ConfirmationDialogObject `json:"confirm,omitempty"`
	ResponseURLEnabled           bool                      `json:"response_url_enabled,omitempty"`
	Filter                       *FilterObject             `json:"filter,omitempty"`
}

// ChannelsSelectBlockElement is a Select menu element with a list of public channels visible to the current user.
// This has a type of "conversations_select."
type ChannelsSelectBlockElement struct {
	blockElement
	Placeholder        *TextCompositionObject    `json:"placeholder"`
	ActionID           ActionID                  `json:"action_id"`
	InitialChannelID   ChannelID                 `json:"initial_channel,omitempty"`
	Confirm            *ConfirmationDialogObject `json:"confirm,omitempty"`
	ResponseURLEnabled bool                      `json:"response_url_enabled,omitempty"`
}
