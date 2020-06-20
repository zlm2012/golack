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
	Verbatim bool   `json:"verbatim, omitempty"`
}

type ConfirmationDialogObject struct {
	compositionObject
	Title   *TextCompositionObject `json:"title"`
	Text    *TextCompositionObject `json:"text"`
	Confirm *TextCompositionObject `json:"confirm"`
	Deny    *TextCompositionObject `json:"deny"`
	Style   string                 `json:"style,omitempty"`
}

type OptionObject struct {
	compositionObject
	Text        *TextCompositionObject `json:"text"`
	Value       string                 `json:"value"`
	Description *TextCompositionObject `json:"description"`
	URL         string                 `json:"url,omitempty"`
}

type OptionGroupObject struct {
	compositionObject
	Label   *TextCompositionObject `json:"label"`
	Options []CompositionObject    `json:"options"`
}

type FilterObject struct {
	compositionObject
	Include                       []string `json:"include,omitempty"` // im, mpim, private or public
	ExcludeExternalSharedChannels bool     `json:"exclude_external_shared_channels,omitempty"`
	ExcludeBotUsers               bool     `json:"exclude_bot_users"`
}
