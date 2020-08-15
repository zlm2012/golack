package webapi

import (
	"github.com/oklahomer/golack/v2/event"
)

// ParseMode defines the parse parameter values for post.message method.
// See https://api.slack.com/docs/message-formatting#parsing_modes
type ParseMode string

const (
	ParseModeNone ParseMode = "none"
	ParseModeFull ParseMode = "full"
)

// String returns a stringified form of BotType
func (mode ParseMode) String() string {
	return string(mode)
}

type AttachmentField struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value"`
	Short bool   `json:"short,omitempty"`
}

type MessageAttachment struct {
	Fallback   string             `json:"fallback"`
	Color      string             `json:"color,omitempty"`
	Pretext    string             `json:"pretext,omitempty"`
	AuthorName string             `json:"author_name,omitempty"`
	AuthorLink string             `json:"author_link,omitempty"`
	AuthorIcon string             `json:"author_icon,omitempty"`
	Title      string             `json:"title,omitempty"`
	TitleLink  string             `json:"title_link,omitempty"`
	Text       string             `json:"text,omitempty"`
	Fields     []*AttachmentField `json:"fields"`
	ImageURL   string             `json:"image_url,omitempty"`
	ThumbURL   string             `json:"thumb_url,omitempty"`
}

// PostMessage is a payload to be sent with chat.postMessage method.
// See https://api.slack.com/methods/chat.postMessage
type PostMessage struct {
	ChannelID       event.ChannelID      `json:"channel"`
	Text            string               `json:"text"`
	AsUser          bool                 `json:"as_user,omitempty"`
	Attachments     []*MessageAttachment `json:"attachments,omitempty"`
	Blocks          []event.Block        `json:"blocks,omitempty"`
	IconEmoji       string               `json:"icon_emoji,omitempty"`
	IconURL         string               `json:"icon_url,omitempty"`
	LinkNames       int                  `json:"link_names,omitempty"`
	Markdown        bool                 `json:"mrkdwn,omitempty"`
	Parse           ParseMode            `json:"parse,omitempty"`
	ReplyBroadcast  bool                 `json:"reply_broadcast,omitempty"`
	ThreadTimeStamp string               `json:"thread_ts,omitempty"`
	UnfurlLinks     bool                 `json:"unfurl_links,omitempty"`
	UnfurlMedia     bool                 `json:"unfurl_media,omitempty"`
	UserName        string               `json:"username,omitempty"`
}

// WithAsUser sets optional boolean value so the outgoing message is sent as a user.
// See https://api.slack.com/methods/chat.postMessage
func (message *PostMessage) WithAsUser(flg bool) *PostMessage {
	message.AsUser = flg
	return message
}

// WithAttachments sets/overrides attachments parameter for current PostMessage.
// See https://api.slack.com/docs/message-attachments
func (message *PostMessage) WithAttachments(attachments []*MessageAttachment) *PostMessage {
	message.Attachments = attachments
	return message
}

// WithBlocks sets/overrides blocks parameter for current PostMessage.
// See https://api.slack.com/messaging/composing/layouts#adding-blocks
func (message *PostMessage) WithBlocks(blocks []event.Block) *PostMessage {
	message.Blocks = blocks
	return message
}

// WithIconEmoji sets an icon emoji for the message.
// See https://api.slack.com/methods/chat.postMessage
func (message *PostMessage) WithIconEmoji(iconEmoji string) *PostMessage {
	message.IconEmoji = iconEmoji
	return message
}

// WithIconURL sets an icon image url for the message.
// See https://api.slack.com/methods/chat.postMessage
func (message *PostMessage) WithIconURL(iconURL string) *PostMessage {
	message.IconURL = iconURL
	return message
}

// WithLinkNames sets/overrides link_names parameter for current PostMessage.
// See https://api.slack.com/methods/chat.postMessage#formatting
func (message *PostMessage) WithLinkNames(linkNames int) *PostMessage {
	message.LinkNames = linkNames
	return message
}

// WithMarkdown sets optional boolean value to decide if the message should be treated as Markdown.
// See https://api.slack.com/methods/chat.postMessage
func (message *PostMessage) WithMarkdown(flg bool) *PostMessage {
	message.Markdown = flg
	return message
}

// WithParse sets/overrides parse parameter for current PostMessage.
// See https://api.slack.com/docs/message-formatting#parsing_modes
func (message *PostMessage) WithParse(parse ParseMode) *PostMessage {
	message.Parse = parse
	return message
}

// WithReplyBroadcast sets optional boolean value so the thread response can be broadcasted.
// Thread identifier must be present with WithThreadTimeStamp() to use this option.
// See https://api.slack.com/docs/message-threading#using_the_web_api
func (message *PostMessage) WithReplyBroadcast(flg bool) *PostMessage {
	message.ReplyBroadcast = flg
	return message
}

// WithThreadTimeStamp sets given ts value to payload.
// See https://api.slack.com/docs/message-threading#using_the_web_api
func (message *PostMessage) WithThreadTimeStamp(ts string) *PostMessage {
	message.ThreadTimeStamp = ts
	return message
}

// WithUnfurlLinks sets/overrides unfurl_links for current PostMessage.
// See https://api.slack.com/docs/message-attachments#unfurling
func (message *PostMessage) WithUnfurlLinks(unfurl bool) *PostMessage {
	message.UnfurlLinks = unfurl
	return message
}

// WithUnfurlLinks sets/overrides unfurl_media for current PostMessage.
// See https://api.slack.com/docs/message-attachments#unfurling
func (message *PostMessage) WithUnfurlMedia(unfurl bool) *PostMessage {
	message.UnfurlMedia = unfurl
	return message
}

// WithUserName sets a name to be used as user name.
// See https://api.slack.com/methods/chat.postMessage
func (message *PostMessage) WithUserName(name string) *PostMessage {
	message.UserName = name
	return message
}

// NewPostMessage creates PostMessage instance with given channel and text settings.
// By default this sets commonly used settings as much as possible. e.g. link_names=1, unfurl_links=true, etc...
// To override those settings and add some extra settings including username, icon_url, or icon_emoji, call setter methods start with With***.
func NewPostMessage(channelID event.ChannelID, text string) *PostMessage {
	return &PostMessage{
		ChannelID:   channelID,
		Text:        text,
		Parse:       ParseModeFull,
		LinkNames:   1,
		UnfurlLinks: true,
		UnfurlMedia: true,
	}
}
