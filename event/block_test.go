package event

import (
	"github.com/oklahomer/golack/v2/testutil"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

var expectedBlock = map[string]interface{}{
	"section": &SectionBlock{
		block: block{
			Type:    "section",
			BlockID: "text1",
		},
		Text: &TextCompositionObject{
			Type:     "mrkdwn",
			Text:     "A message *with some bold text* and _some italicized text_.",
			Emoji:    false,
			Verbatim: false,
		},
		Fields:    nil,
		Accessory: nil,
	},
	"section.text": &SectionBlock{
		block: block{
			Type: "section",
		},
		Text: &TextCompositionObject{
			Type: "mrkdwn",
			Text: "A message *with some bold text* and _some italicized text_.",
		},
		Fields: []*TextCompositionObject{
			{
				Type: "mrkdwn",
				Text: "High",
			},
			{
				Type:  "plain_text",
				Emoji: true,
				Text:  "String",
			},
		},
		Accessory: nil,
	},
	"section.datepicker": &SectionBlock{
		block: block{
			Type: "section",
		},
		Text: &TextCompositionObject{
			Type: "mrkdwn",
			Text: "*Sally* has requested you set the deadline for the Nano launch project",
		},
		Accessory: &DatePickerBlockElement{
			blockElement: blockElement{
				Type: "datepicker",
			},
			ActionID:    "datepicker123",
			InitialDate: "1990-04-28",
			PlaceHolder: &TextCompositionObject{
				Type: "plain_text",
				Text: "Select a date",
			},
		},
	},
	"divider": &DividerBlock{
		block: block{
			Type: "divider",
		},
	},
	"image": &ImageBlock{
		block: block{
			Type:    "image",
			BlockID: "image4",
		},
		ImageURL: "http://placekitten.com/500/500",
		AltText:  "An incredibly cute kitten.",
		Title: &TextCompositionObject{
			Type: "plain_text",
			Text: "Please enjoy this photo of a kitten",
		},
	},
	"actions": &ActionsBlock{
		block: block{
			Type:    "actions",
			BlockID: "actions1",
		},
		Elements: []BlockElement{
			&StaticSelectBlockElement{
				blockElement: blockElement{
					Type: "static_select",
				},
				Placeholder: &TextCompositionObject{
					Type: "plain_text",
					Text: "Which witch is the witchiest witch?",
				},
				ActionID: "select_2",
				Options: []*OptionObject{
					{
						Text: &TextCompositionObject{
							Type: "plain_text",
							Text: "Matilda",
						},
						Value: "matilda",
					},
					{
						Text: &TextCompositionObject{
							Type: "plain_text",
							Text: "Glinda",
						},
						Value: "glinda",
					},
					{
						Text: &TextCompositionObject{
							Type: "plain_text",
							Text: "Granny Weatherwax",
						},
						Value: "grannyWeatherwax",
					},
					{
						Text: &TextCompositionObject{
							Type: "plain_text",
							Text: "Hermione",
						},
						Value: "hermione",
					},
				},
			},
			&ButtonBlockElement{
				blockElement: blockElement{
					Type: "button",
				},
				Text: &TextCompositionObject{
					Type: "plain_text",
					Text: "Cancel",
				},
				Value:    "cancel",
				ActionID: "button_1",
			},
		},
	},
	"context": &ContextBlock{
		block: block{
			Type: "context",
		},
		Elements: []BlockElement{
			&ImageBlockElement{
				blockElement: blockElement{
					Type: "image",
				},
				ImageURL: "https://image.freepik.com/free-photo/red-drawing-pin_1156-445.jpg",
				AltText:  "images",
			},
			&TextObjectBlockElement{
				blockElement: blockElement{
					Type: "mrkdwn",
				},
				Text: "Location: **Dogpatch**",
			},
		},
	},
	"file": &FileBlock{
		block: block{
			Type: "file",
		},
		ExternalID: "ABCD1",
		Source:     "remote",
	},
}

func TestUnmarshalBlock(t *testing.T) {
	directory := filepath.Join("..", "testdata", "event", "block")

	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Fatalf("Failed to read testdata directory: %s", err.Error())
		}

		if !strings.HasSuffix(path, ".json.golden") {
			// Skip irrelevant file
			return nil
		}

		checked := map[string]bool{}
		filename := filepath.Base(path)
		t.Run(filename, func(t *testing.T) {
			input, err := ioutil.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to read file: %s. Error: %s.", path, err.Error())
			}

			blockType := strings.TrimSuffix(filename, ".json.golden")
			checked[blockType] = true

			decoded, err := UnmarshalBlock(input)
			if err != nil {
				t.Fatalf("Failed to decode JSON: %s. Error: %s.", path, err.Error())
			}

			expected, ok := expectedBlock[blockType]
			if !ok {
				t.Fatalf("Expected payload for %s is not defined in test", blockType)
			}

			testutil.Compare([]string{blockType}, reflect.ValueOf(expected), reflect.ValueOf(decoded), t)
		})

		return nil
	})

}
