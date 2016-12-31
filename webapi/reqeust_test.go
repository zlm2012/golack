package webapi

import (
	"strconv"
	"testing"
)

func TestParseMode_String(t *testing.T) {
	testVars := []struct {
		mode ParseMode
		val  string
	}{
		{
			mode: ParseModeNone,
			val:  "none",
		},
		{
			mode: ParseModeFull,
			val:  "full",
		},
	}

	for i, testVar := range testVars {
		if testVar.mode.String() != testVar.val {
			t.Errorf("expected value is not returned on test #%d: %s", i+1, testVar.mode.String())
		}
	}
}

func TestNewPostMessage(t *testing.T) {
	channel := "myChannel"
	text := "myText"
	message := NewPostMessage(channel, text)

	if message == nil {
		t.Fatal("new PostMessage instance is not returned.")
	}

	if message.Channel != channel {
		t.Errorf("expected channel is not set: %s.", message.Channel)
	}

	if message.Text != text {
		t.Errorf("expected text is not set: %s.", message.Text)
	}
}

func TestNewPostMessageWithAttachments(t *testing.T) {
	channel := "myChannel"
	text := "myText"
	attachment := &MessageAttachment{}
	message := NewPostMessageWithAttachments(channel, text, []*MessageAttachment{attachment})

	if message == nil {
		t.Fatal("new PostMessage instance is not returned.")
	}

	if message.Channel != channel {
		t.Errorf("expected channel is not set: %s.", message.Channel)
	}

	if message.Text != text {
		t.Errorf("expected text is not set: %s.", message.Text)
	}

	if message.Attachments == nil || len(message.Attachments) == 0 {
		t.Fatal("given MessageAttachment is not set.")
	}

	if message.Attachments[0] != attachment {
		t.Errorf("expected attachment is not set: %#v.", message.Attachments[0])
	}
}

func TestPostMessage_WithLinkNames(t *testing.T) {
	message := &PostMessage{}
	oldVal := message.LinkNames
	message.WithLinkNames(123)

	if message.LinkNames == oldVal {
		t.Error("value is not updated.")
	}
}

func TestPostMessage_WithParse(t *testing.T) {
	message := &PostMessage{}
	message.WithParse(ParseModeFull)

	if message.Parse == "" {
		t.Error("value is not updated.")
	}
}

func TestPostMessage_WithUnfurlLinks(t *testing.T) {
	message := &PostMessage{}
	oldVal := message.UnfurlLinks
	message.WithUnfurlLinks(!oldVal)

	if message.UnfurlLinks == oldVal {
		t.Error("value is not updated.")
	}
}

func TestPostMessage_WithUnfurlMedia(t *testing.T) {
	message := &PostMessage{}
	oldVal := message.UnfurlMedia
	message.WithUnfurlMedia(!oldVal)

	if message.UnfurlMedia == oldVal {
		t.Error("value is not updated.")
	}
}

func TestPostMessage_ToURLValues(t *testing.T) {
	channel := "myChannel"
	text := "myText"
	parse := ParseModeFull
	linkNames := 1
	attachment := &MessageAttachment{}
	unfurlLinks := true
	unfurlMedia := true
	userName := "myName"
	asUser := true
	iconUrl := "http://example.com/icon.png"
	iconEmoji := ":chart_with_upwards_trend:"
	message := &PostMessage{
		Channel:     channel,
		Text:        text,
		Parse:       parse,
		LinkNames:   linkNames,
		Attachments: []*MessageAttachment{attachment},
		UnfurlLinks: unfurlLinks,
		UnfurlMedia: unfurlMedia,
		UserName:    userName,
		AsUser:      asUser,
		IconURL:     iconUrl,
		IconEmoji:   iconEmoji,
	}

	testVars := []struct {
		key string
		val interface{}
	}{
		{
			key: "channel",
			val: channel,
		},
		{
			key: "text",
			val: text,
		},
		{
			key: "parse",
			val: parse.String(),
		},
		{
			key: "link_names",
			val: strconv.Itoa(linkNames),
		},
		{
			key: "unfurl_links",
			val: strconv.FormatBool(unfurlLinks),
		},
		{
			key: "unfurl_media",
			val: strconv.FormatBool(unfurlMedia),
		},
		{
			key: "username",
			val: userName,
		},
		{
			key: "as_user",
			val: strconv.FormatBool(asUser),
		},
		{
			key: "icon_url",
			val: iconUrl,
		},
		{
			key: "icon_emoji",
			val: iconEmoji,
		},
	}

	urlVal := message.ToURLValues()
	for _, testVar := range testVars {
		if urlVal.Get(testVar.key) != testVar.val {
			t.Errorf("expected value is not returned. key: %s. val: %#v.", testVar.key, urlVal.Get(testVar.key))
		}
	}
	// TODO check marshaled attachments field
}
