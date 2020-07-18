package event

// List of composition objects: https://api.slack.com/reference/block-kit/composition-objects

type CompositionObject interface {
	composition()
}

type compositionObject struct {
}

var _ CompositionObject = (*compositionObject)(nil)

func (co *compositionObject) composition() {}

type TextCompositionObject struct {
	compositionObject
	Type     string `json:"type"`
	Text     string `json:"text"`
	Emoji    bool   `json:"emoji,omitempty"`
	Verbatim bool   `json:"verbatim,omitempty"`
}

func (tc *TextCompositionObject) WithEmoji(flg bool) *TextCompositionObject {
	tc.Emoji = flg
	return tc
}

func (tc *TextCompositionObject) WithVerbatim(flg bool) *TextCompositionObject {
	tc.Verbatim = flg
	return tc
}

func NewPlainTextCompositionObject(text string) *TextCompositionObject {
	return &TextCompositionObject{
		Type:     "plain_text",
		Text:     text,
		Emoji:    false,
		Verbatim: false,
	}
}

func NewMarkdownTextCompositionObject(text string) *TextCompositionObject {
	return &TextCompositionObject{
		Type: "mrkdwn",
		Text: text,
	}
}

type ConfirmationDialogObject struct {
	compositionObject
	Title   *TextCompositionObject `json:"title"`
	Text    *TextCompositionObject `json:"text"`
	Confirm *TextCompositionObject `json:"confirm"`
	Deny    *TextCompositionObject `json:"deny"`
	Style   Style                  `json:"style,omitempty"`
}

func (cd *ConfirmationDialogObject) WithStyle(style Style) *ConfirmationDialogObject {
	cd.Style = style
	return cd
}

func NewConfirmationDialogObject(title, text, confirm, deny *TextCompositionObject) *ConfirmationDialogObject {
	return &ConfirmationDialogObject{
		Title:   title,
		Text:    text,
		Confirm: confirm,
		Deny:    deny,
	}
}

type OptionObject struct {
	compositionObject
	Text        *TextCompositionObject `json:"text"`
	Value       string                 `json:"value"`
	Description *TextCompositionObject `json:"description,omitempty"`
	URL         string                 `json:"url,omitempty"`
}

func (o *OptionObject) WithDescription(description *TextCompositionObject) *OptionObject {
	o.Description = description
	return o
}

func (o *OptionObject) WithURL(url string) *OptionObject {
	o.URL = url
	return o
}

func NewOptionObject(text *TextCompositionObject, value string) *OptionObject {
	return &OptionObject{
		Text:  text,
		Value: value,
	}
}

type OptionGroupObject struct {
	compositionObject
	Label   *TextCompositionObject `json:"label"`
	Options []CompositionObject    `json:"options"`
}

func NewOptionGroupObject(label *TextCompositionObject, options []CompositionObject) *OptionGroupObject {
	return &OptionGroupObject{
		Label:   label,
		Options: options,
	}
}

type FilterObject struct {
	compositionObject
	Include                       []ConversationType `json:"include,omitempty"`
	ExcludeExternalSharedChannels bool               `json:"exclude_external_shared_channels,omitempty"`
	ExcludeBotUsers               bool               `json:"exclude_bot_users,omitempty"`
}

func (f *FilterObject) WithInclude(conversationTypes []ConversationType) *FilterObject {
	f.Include = conversationTypes
	return f
}

func (f *FilterObject) WithExcludeExternalSharedChannels(flg bool) *FilterObject {
	f.ExcludeExternalSharedChannels = flg
	return f
}

func (f *FilterObject) WithExcludeBotUsers(flg bool) *FilterObject {
	f.ExcludeBotUsers = flg
	return f
}

func NewFilterObject() *FilterObject {
	return &FilterObject{
		Include:                       nil,
		ExcludeExternalSharedChannels: false,
		ExcludeBotUsers:               false,
	}
}
