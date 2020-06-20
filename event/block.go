package event

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"golang.org/x/xerrors"
)

// List of layout blocks: https://api.slack.com/reference/block-kit/blocks

func UnmarshalBlock(input json.RawMessage) (Block, error) {
	parsed := gjson.ParseBytes(input)

	typeValue := parsed.Get("type")

	var typed Block
	switch t := typeValue.String(); t {
	case "section":
		typed = &SectionBlock{}

	case "divider":
		typed = &DividerBlock{}
	case "image":
		typed = &ImageBlock{}

	case "actions":
		typed = &ActionBlock{}

	case "context":
		typed = &ContextBlock{}

	case "input":
		typed = &InputBlock{}

	case "file":
		typed = &FileBlock{}

	default:
		return nil, fmt.Errorf("failed to handle unknown block type: %s", t)
	}

	err := json.Unmarshal(input, typed)
	if err != nil {
		return nil, xerrors.Errorf("failed to unmarshal %T", typed)
	}
	return typed, nil
}

type Block interface {
	BlockType() string
}

type block struct {
	Type    string  `json:"type"`
	BlockID BlockID `json:"block_id,omitempty"`
}

var _ Block = (*block)(nil)

func (b *block) BlockType() string {
	return b.Type
}

type ActionBlock struct {
	block
	Elements []BlockElement `json:"elements"`
}

func (ab *ActionBlock) UnmarshalJSON(b []byte) error {
	type alias ActionBlock
	t := &struct {
		*alias
		Elements []json.RawMessage `json:"elements"`
	}{
		alias: (*alias)(ab),
	}
	err := json.Unmarshal(b, t)
	if err != nil {
		return err
	}

	var elements []BlockElement
	for _, elem := range t.Elements {
		element, err := UnmarshalBlockElement(elem)
		if err != nil {
			return xerrors.Errorf("failed to unmarshal given element: %w", err)
		}
		elements = append(elements, element)
	}

	ab.Elements = elements
	return nil
}

type ContextBlock struct {
	block
	Elements []BlockElement `json:"elements"`
}

func (cb *ContextBlock) UnmarshalJSON(b []byte) error {
	type alias ContextBlock
	t := &struct {
		*alias
		Elements []json.RawMessage `json:"elements"`
	}{
		alias: (*alias)(cb),
	}
	err := json.Unmarshal(b, t)
	if err != nil {
		return err
	}

	var elements []BlockElement
	for _, elem := range t.Elements {
		element, err := UnmarshalBlockElement(elem)
		if err != nil {
			return xerrors.Errorf("failed to unmarshal given element: %w", err)
		}
		elements = append(elements, element)
	}

	cb.Elements = elements
	return nil
}

type DividerBlock struct {
	block
}

type FileBlock struct {
	block
	ExternalID string `json:"external_id"`
	Source     string `json:"source"`
}

type ImageBlock struct {
	block
	ImageURL string                 `json:"image_url"`
	AltText  string                 `json:"alt_text"`
	Title    *TextCompositionObject `json:"title,omitempty"`
}

type InputBlock struct {
	block
	Label    *TextCompositionObject `json:"label"`
	Element  BlockElement           `json:"element"`
	Hint     *TextCompositionObject `json:"hint,omitempty"`
	Optional bool                   `json:"optional,omitempty"`
}

func (ib *InputBlock) UnmarshalJSON(b []byte) error {
	type alias InputBlock
	t := &struct {
		*alias
		Element json.RawMessage `json:"element"`
	}{
		alias: (*alias)(ib),
	}
	err := json.Unmarshal(b, t)
	if err != nil {
		return err
	}

	element, err := UnmarshalBlockElement(t.Element)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal given element: %w", err)
	}
	ib.Element = element
	return nil
}

type SectionBlock struct {
	block
	Text      *TextCompositionObject   `json:"text"`
	Fields    []*TextCompositionObject `json:"fields"`
	Accessory CompositionObject        `json:"accessory"`
}
