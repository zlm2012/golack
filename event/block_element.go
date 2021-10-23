package event

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
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
		typed = &MultiStaticSelectBlockElement{}

	case "multi_external_select":
		typed = &MultiExternalSelectBlockElement{}

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

	case "mrkdwn", "plain_text":
		typed = &TextObjectBlockElement{}

	default:
		return nil, fmt.Errorf("failed to handle unknown block element type: %s", t)
	}

	err := json.Unmarshal(input, typed)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal %T", typed)
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
	Style    Style                     `json:"style,omitempty"`
	Confirm  *ConfirmationDialogObject `json:"confirm,omitempty"`
}

func (b *ButtonBlockElement) WithURL(url string) *ButtonBlockElement {
	b.URL = url
	return b
}

func (b *ButtonBlockElement) WithValue(value string) *ButtonBlockElement {
	b.Value = value
	return b
}

func (b *ButtonBlockElement) WithStyle(style Style) *ButtonBlockElement {
	b.Style = style
	return b
}

func (b *ButtonBlockElement) WithConfirm(confirm *ConfirmationDialogObject) *ButtonBlockElement {
	b.Confirm = confirm
	return b
}

func NewButtonBlockElement(text *TextCompositionObject, actionID ActionID) *ButtonBlockElement {
	return &ButtonBlockElement{
		blockElement: blockElement{
			Type: "button",
		},
		Text:     text,
		ActionID: actionID,
	}
}

type CheckboxBlockElement struct {
	blockElement
	ActionID       ActionID                  `json:"action_id"`
	Options        []*OptionObject           `json:"options"`
	InitialOptions []*OptionObject           `json:"initial_options,omitempty"`
	Confirm        *ConfirmationDialogObject `json:"confirm,omitempty"`
}

func (c *CheckboxBlockElement) WithInitialOptions(options []*OptionObject) *CheckboxBlockElement {
	c.InitialOptions = options
	return c
}

func (c *CheckboxBlockElement) WithConfirm(confirm *ConfirmationDialogObject) *CheckboxBlockElement {
	c.Confirm = confirm
	return c
}

func NewCheckboxBlockElement(actionID ActionID, options []*OptionObject) *CheckboxBlockElement {
	return &CheckboxBlockElement{
		blockElement: blockElement{
			Type: "checkboxes",
		},
		ActionID: actionID,
		Options:  options,
	}
}

type DatePickerBlockElement struct {
	blockElement
	ActionID    ActionID                  `json:"action_id"`
	PlaceHolder *TextCompositionObject    `json:"placeholder,omitempty"`
	InitialDate string                    `json:"initial_date,omitempty"`
	Confirm     *ConfirmationDialogObject `json:"confirm,omitempty"`
}

func (d *DatePickerBlockElement) WithPlaceHolder(placeholder *TextCompositionObject) *DatePickerBlockElement {
	d.PlaceHolder = placeholder
	return d
}

func (d *DatePickerBlockElement) WithInitialDate(date string) *DatePickerBlockElement {
	d.InitialDate = date
	return d
}

func (d *DatePickerBlockElement) WithConfirm(confirm *ConfirmationDialogObject) *DatePickerBlockElement {
	d.Confirm = confirm
	return d
}

func NewDatePickerBlockElement(actionID ActionID) *DatePickerBlockElement {
	return &DatePickerBlockElement{
		blockElement: blockElement{
			Type: "datepicker",
		},
		ActionID: actionID,
	}
}

type ImageBlockElement struct {
	blockElement
	ImageURL string `json:"image_url"`
	AltText  string `json:"alt_text"`
}

func NewImageBlockElement(imageURL string, altText string) *ImageBlockElement {
	return &ImageBlockElement{
		blockElement: blockElement{
			Type: "image",
		},
		ImageURL: imageURL,
		AltText:  altText,
	}
}

// MultiStaticSelectBlockElement is a Multi-select menu with static options.
// This has a type of "multi_static_select."
type MultiStaticSelectBlockElement struct {
	blockElement
	Placeholder      *TextCompositionObject    `json:"placeholder"`
	ActionID         ActionID                  `json:"action_id"`
	Options          []*OptionObject           `json:"options"`
	OptionGroups     []*OptionGroupObject      `json:"option_groups,omitempty"`
	InitialOptions   []*OptionObject           `json:"initial_options,omitempty"`
	Confirm          *ConfirmationDialogObject `json:"confirm,omitempty"`
	MaxSelectedItems int                       `json:"max_selected_items,omitempty"`
}

func (m *MultiStaticSelectBlockElement) WithOptionGroups(optionGroups []*OptionGroupObject) *MultiStaticSelectBlockElement {
	m.OptionGroups = optionGroups
	return m
}

func (m *MultiStaticSelectBlockElement) WithInitialOptions(options []*OptionObject) *MultiStaticSelectBlockElement {
	m.InitialOptions = options
	return m
}

func (m *MultiStaticSelectBlockElement) WithConfirm(confirm *ConfirmationDialogObject) *MultiStaticSelectBlockElement {
	m.Confirm = confirm
	return m
}

func (m *MultiStaticSelectBlockElement) WithMaxSelectedItems(max int) *MultiStaticSelectBlockElement {
	m.MaxSelectedItems = max
	return m
}

func NewMultiStaticSelectBlockElement(placeholder *TextCompositionObject, actionID ActionID, options []*OptionObject) *MultiStaticSelectBlockElement {
	return &MultiStaticSelectBlockElement{
		blockElement: blockElement{
			Type: "multi_static_select",
		},
		Placeholder: placeholder,
		ActionID:    actionID,
		Options:     options,
	}
}

// MultiExternalSelectBlockElement is a Multi-select menu with external data source.
// This has a type of "multi_external_select."
type MultiExternalSelectBlockElement struct {
	blockElement
	Placeholder      *TextCompositionObject    `json:"placeholder"`
	ActionID         ActionID                  `json:"action_id"`
	MinQueryLength   int                       `json:"min_query_length,omitempty"`
	InitialOptions   []*OptionObject           `json:"initial_options,omitempty"`
	Confirm          *ConfirmationDialogObject `json:"confirm,omitempty"`
	MaxSelectedItems int                       `json:"max_selected_items,omitempty"`
}

func (m *MultiExternalSelectBlockElement) WithMinQueryLength(min int) *MultiExternalSelectBlockElement {
	m.MinQueryLength = min
	return m
}

func (m *MultiExternalSelectBlockElement) WithInitialOptions(options []*OptionObject) *MultiExternalSelectBlockElement {
	m.InitialOptions = options
	return m
}

func (m *MultiExternalSelectBlockElement) WithConfirm(confirm *ConfirmationDialogObject) *MultiExternalSelectBlockElement {
	m.Confirm = confirm
	return m
}

func (m *MultiExternalSelectBlockElement) WithMaxSelectedItems(max int) *MultiExternalSelectBlockElement {
	m.MaxSelectedItems = max
	return m
}

func NewMultiExternalSelectBlockElement(placeholder *TextCompositionObject, actionID ActionID) *MultiExternalSelectBlockElement {
	return &MultiExternalSelectBlockElement{
		blockElement: blockElement{
			Type: "multi_external_select",
		},
		Placeholder: placeholder,
		ActionID:    actionID,
	}
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

func (m *MultiUsersSelectBlockElement) WithInitialUserIDs(ids []UserID) *MultiUsersSelectBlockElement {
	m.InitialUserIDs = ids
	return m
}

func (m *MultiUsersSelectBlockElement) WithConfirm(confirm *ConfirmationDialogObject) *MultiUsersSelectBlockElement {
	m.Confirm = confirm
	return m
}

func (m *MultiUsersSelectBlockElement) WithMaxSelectedItems(max int) *MultiUsersSelectBlockElement {
	m.MaxSelectedItems = max
	return m
}

func NewMultiUsersSelectBlockElement(placeholder *TextCompositionObject, actionID ActionID) *MultiUsersSelectBlockElement {
	return &MultiUsersSelectBlockElement{
		blockElement: blockElement{
			Type: "multi_users_select",
		},
		Placeholder: placeholder,
		ActionID:    actionID,
	}
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

func (m *MultiConversationsSelectBlockElement) WithInitialConversations(conversations []string) *MultiConversationsSelectBlockElement {
	m.InitialConversations = conversations
	return m
}

func (m *MultiConversationsSelectBlockElement) WithDefaultToCurrentConversation(flg bool) *MultiConversationsSelectBlockElement {
	m.DefaultToCurrentConversation = flg
	return m
}

func (m *MultiConversationsSelectBlockElement) WithConfirm(confirm *ConfirmationDialogObject) *MultiConversationsSelectBlockElement {
	m.Confirm = confirm
	return m
}

func (m *MultiConversationsSelectBlockElement) WithMaxSelectedItems(max int) *MultiConversationsSelectBlockElement {
	m.MaxSelectedItems = max
	return m
}

func (m *MultiConversationsSelectBlockElement) WithFilter(filter *FilterObject) *MultiConversationsSelectBlockElement {
	m.Filter = filter
	return m
}

func NewMultiConversationsSelectBlockElement(placeholder *TextCompositionObject, actionID ActionID) *MultiConversationsSelectBlockElement {
	return &MultiConversationsSelectBlockElement{
		blockElement: blockElement{
			Type: "multi_conversations_select",
		},
		Placeholder: placeholder,
		ActionID:    actionID,
	}
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

func (m *MultiChannelsSelectBlockElement) WithInitialChannelIDs(ids []ChannelID) *MultiChannelsSelectBlockElement {
	m.InitialChannelIDs = ids
	return m
}

func (m *MultiChannelsSelectBlockElement) WithConfirm(confirm *ConfirmationDialogObject) *MultiChannelsSelectBlockElement {
	m.Confirm = confirm
	return m
}

func (m *MultiChannelsSelectBlockElement) WithMaxSelectedItems(max int) *MultiChannelsSelectBlockElement {
	m.MaxSelectedItems = max
	return m
}

func NewMultiChannelsSelectBlockElement(placeholder *TextCompositionObject, actionID ActionID) *MultiChannelsSelectBlockElement {
	return &MultiChannelsSelectBlockElement{
		blockElement: blockElement{
			Type: "multi_channels_select",
		},
		Placeholder: placeholder,
		ActionID:    actionID,
	}
}

type OverflowBlockElement struct {
	blockElement
	ActionID ActionID                  `json:"action_id"`
	Options  []*OptionObject           `json:"options"`
	Confirm  *ConfirmationDialogObject `json:"confirm"`
}

func NewOverflowBlockElement(actionID ActionID, options []*OptionObject, confirm *ConfirmationDialogObject) *OverflowBlockElement {
	return &OverflowBlockElement{
		blockElement: blockElement{
			Type: "overflow",
		},
		ActionID: actionID,
		Options:  options,
		Confirm:  confirm,
	}
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

func (p *PlainTextInputBlockElement) WithPlaceholder(placeholder *TextCompositionObject) *PlainTextInputBlockElement {
	p.Placeholder = placeholder
	return p
}

func (p *PlainTextInputBlockElement) WithInitialValue(initialValue string) *PlainTextInputBlockElement {
	p.InitialValue = initialValue
	return p
}

func (p *PlainTextInputBlockElement) WithMultiline(flg bool) *PlainTextInputBlockElement {
	p.Multiline = flg
	return p
}

func (p *PlainTextInputBlockElement) WithMinLength(min int) *PlainTextInputBlockElement {
	p.MinLength = min
	return p
}

func (p *PlainTextInputBlockElement) WithMaxLength(max int) *PlainTextInputBlockElement {
	p.MaxLength = max
	return p
}

func NewPlainTextInputBlockElement(actionID ActionID) *PlainTextInputBlockElement {
	return &PlainTextInputBlockElement{
		blockElement: blockElement{
			Type: "plain_text_input",
		},
		ActionID: actionID,
	}
}

type RadioButtonGroupBlockElement struct {
	blockElement
	ActionID      ActionID                  `json:"action_id"`
	Options       []*OptionObject           `json:"options"`
	InitialOption *OptionObject             `json:"initial_option,omitempty"`
	Confirm       *ConfirmationDialogObject `json:"confirm,omitempty"`
}

func (r *RadioButtonGroupBlockElement) WithInitialOption(option *OptionObject) *RadioButtonGroupBlockElement {
	r.InitialOption = option
	return r
}

func (r *RadioButtonGroupBlockElement) WithConfirm(confirm *ConfirmationDialogObject) *RadioButtonGroupBlockElement {
	r.Confirm = confirm
	return r
}

func NewRadioButtonGroupBlockElement(actionID ActionID, options []*OptionObject) *RadioButtonGroupBlockElement {
	return &RadioButtonGroupBlockElement{
		blockElement: blockElement{
			Type: "radio_buttons",
		},
		ActionID: actionID,
		Options:  options,
	}
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

func (s *StaticSelectBlockElement) WithOptionGroups(optionGroups []*OptionGroupObject) *StaticSelectBlockElement {
	s.OptionGroups = optionGroups
	return s
}

func (s *StaticSelectBlockElement) WithInitialOption(option *OptionObject) *StaticSelectBlockElement {
	s.InitialOption = option
	return s
}

func (s *StaticSelectBlockElement) WithConfirm(confirm *ConfirmationDialogObject) *StaticSelectBlockElement {
	s.Confirm = confirm
	return s
}

func NewStaticSelectBlockElement(placeholder *TextCompositionObject, actionID ActionID, options []*OptionObject) *StaticSelectBlockElement {
	return &StaticSelectBlockElement{
		blockElement: blockElement{
			Type: "static_select",
		},
		Placeholder: placeholder,
		ActionID:    actionID,
		Options:     options,
	}
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

func (e *ExternalSelectBlockElement) WithInitialOption(option *OptionObject) *ExternalSelectBlockElement {
	e.InitialOption = option
	return e
}

func (e *ExternalSelectBlockElement) WithMinQueryLength(min int) *ExternalSelectBlockElement {
	e.MinQueryLength = min
	return e
}

func (e *ExternalSelectBlockElement) WithConfirm(confirm *ConfirmationDialogObject) *ExternalSelectBlockElement {
	e.Confirm = confirm
	return e
}

func NewExternalSelectBlockElement(placeholder *TextCompositionObject, actionID ActionID) *ExternalSelectBlockElement {
	return &ExternalSelectBlockElement{
		blockElement: blockElement{
			Type: "external_select",
		},
		Placeholder: placeholder,
		ActionID:    actionID,
	}
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

func (u *UsersSelectBlockElement) WithInitialUserID(userID UserID) *UsersSelectBlockElement {
	u.InitialUserID = userID
	return u
}

func (u *UsersSelectBlockElement) WithConfirm(confirm *ConfirmationDialogObject) *UsersSelectBlockElement {
	u.Confirm = confirm
	return u
}

func NewUsersSelectBlockElement(placeholder *TextCompositionObject, actionID ActionID) *UsersSelectBlockElement {
	return &UsersSelectBlockElement{
		blockElement: blockElement{
			Type: "users_select",
		},
		Placeholder: placeholder,
		ActionID:    actionID,
	}
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

func (c *ConversationsSelectBlockElement) WithInitialConversation(initialConversation string) *ConversationsSelectBlockElement {
	c.InitialConversation = initialConversation
	return c
}

func (c *ConversationsSelectBlockElement) WithDefaultToCurrentConversation(flg bool) *ConversationsSelectBlockElement {
	c.DefaultToCurrentConversation = flg
	return c
}

func (c *ConversationsSelectBlockElement) WithConfirm(confirm *ConfirmationDialogObject) *ConversationsSelectBlockElement {
	c.Confirm = confirm
	return c
}

func (c *ConversationsSelectBlockElement) WithResponseURLEnabled(flg bool) *ConversationsSelectBlockElement {
	c.ResponseURLEnabled = flg
	return c
}

func (c *ConversationsSelectBlockElement) WithFilter(filter *FilterObject) *ConversationsSelectBlockElement {
	c.Filter = filter
	return c
}

func NewConversationsSelectBlockElement(placeholder *TextCompositionObject, actionID ActionID) *ConversationsSelectBlockElement {
	return &ConversationsSelectBlockElement{
		blockElement: blockElement{
			Type: "conversations_select",
		},
		Placeholder: placeholder,
		ActionID:    actionID,
	}
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

func (c *ChannelsSelectBlockElement) WithInitialChannelID(channelID ChannelID) *ChannelsSelectBlockElement {
	c.InitialChannelID = channelID
	return c
}

func (c *ChannelsSelectBlockElement) WithConfirm(confirm *ConfirmationDialogObject) *ChannelsSelectBlockElement {
	c.Confirm = confirm
	return c
}

func (c *ChannelsSelectBlockElement) WithResponseURLEnabled(flg bool) *ChannelsSelectBlockElement {
	c.ResponseURLEnabled = flg
	return c
}

func NewChannelsSelectBlockElement(placeholder *TextCompositionObject, actionID ActionID) *ChannelsSelectBlockElement {
	return &ChannelsSelectBlockElement{
		blockElement: blockElement{
			Type: "channels_select",
		},
		Placeholder: placeholder,
		ActionID:    actionID,
	}
}

// TextObjectBlockElement has equal fields as TextCompositionObject does.
// ContextBlock's Elements can include BlockElement and TextCompositionObject so this is a compromise to let []BlockElement include TextCompositionObject.
type TextObjectBlockElement struct {
	blockElement
	Text     string `json:"text"`
	Emoji    bool   `json:"emoji,omitempty"`
	Verbatim bool   `json:"verbatim,omitempty"`
}

func (t *TextObjectBlockElement) WithEmoji(flg bool) *TextObjectBlockElement {
	t.Emoji = flg
	return t
}

func (t *TextObjectBlockElement) WithVerbatim(flg bool) *TextObjectBlockElement {
	t.Verbatim = flg
	return t
}

func NewPlainTextObjectBlockElement(text string) *TextObjectBlockElement {
	return &TextObjectBlockElement{
		blockElement: blockElement{
			Type: "plain_text",
		},
		Text: text,
	}
}

func NewMarkdownTextObjectBlockElement(text string) *TextObjectBlockElement {
	return &TextObjectBlockElement{
		blockElement: blockElement{
			Type: "mrkdwn",
		},
		Text: text,
	}
}
