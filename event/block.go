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
		typed = &ActionsBlock{}

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
		return nil, xerrors.Errorf("failed to unmarshal %T: %w", typed, err)
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

func (b *block) WithBlockID(blockID BlockID) *block {
	b.BlockID = blockID
	return b
}

type ActionsBlock struct {
	block
	Elements []BlockElement `json:"elements"`
}

func NewActionsBlock(elements []BlockElement) *ActionsBlock {
	return &ActionsBlock{
		block: block{
			Type: "actions",
		},
		Elements: elements,
	}
}

func (ab *ActionsBlock) UnmarshalJSON(b []byte) error {
	type alias ActionsBlock
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

func NewContextBlock(elements []BlockElement) *ContextBlock {
	return &ContextBlock{
		block: block{
			Type: "context",
		},
		Elements: elements,
	}
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

func NewDividerBlock() *DividerBlock {
	return &DividerBlock{
		block: block{
			Type: "divider",
		},
	}
}

type FileBlock struct {
	block
	ExternalID string `json:"external_id"`
	Source     string `json:"source"`
}

func NewRemoteFileBlock(externalID string) *FileBlock {
	return &FileBlock{
		block: block{
			Type: "file",
		},
		ExternalID: externalID,
		Source:     "remote",
	}
}

type ImageBlock struct {
	block
	ImageURL string                 `json:"image_url"`
	AltText  string                 `json:"alt_text"`
	Title    *TextCompositionObject `json:"title,omitempty"`
}

func (ib *ImageBlock) WithTitle(title *TextCompositionObject) *ImageBlock {
	ib.Title = title
	return ib
}

func NewImageBlock(imageURL string, altText string) *ImageBlock {
	return &ImageBlock{
		block: block{
			Type: "image",
		},
		ImageURL: imageURL,
		AltText:  altText,
	}
}

type InputBlock struct {
	block
	Label    *TextCompositionObject `json:"label"`
	Element  BlockElement           `json:"element"`
	Hint     *TextCompositionObject `json:"hint,omitempty"`
	Optional bool                   `json:"optional,omitempty"`
}

func (ib *InputBlock) WithHint(hint *TextCompositionObject) *InputBlock {
	ib.Hint = hint
	return ib
}

func (ib *InputBlock) WithOptional(flg bool) *InputBlock {
	ib.Optional = flg
	return ib
}

func NewInputBlock(label *TextCompositionObject, element BlockElement) *InputBlock {
	return &InputBlock{
		block: block{
			Type: "input",
		},
		Label:   label,
		Element: element,
	}
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

	if len(t.Element) > 0 {
		element, err := UnmarshalBlockElement(t.Element)
		if err != nil {
			return xerrors.Errorf("failed to unmarshal given element: %w", err)
		}
		ib.Element = element
	}

	return nil
}

type SectionBlock struct {
	block
	Text      *TextCompositionObject   `json:"text"`
	Fields    []*TextCompositionObject `json:"fields,omitempty"`
	Accessory BlockElement             `json:"accessory,omitempty"`
}

func (sb *SectionBlock) WithFields(fields []*TextCompositionObject) *SectionBlock {
	sb.Fields = fields
	return sb
}

func (sb *SectionBlock) WithAccessory(accessory BlockElement) *SectionBlock {
	sb.Accessory = accessory
	return sb
}

func NewSectionBlock(text *TextCompositionObject) *SectionBlock {
	return &SectionBlock{
		block: block{
			Type: "section",
		},
		Text: text,
	}
}

func (sb *SectionBlock) UnmarshalJSON(b []byte) error {
	type alias SectionBlock
	t := &struct {
		*alias
		Accessory json.RawMessage `json:"accessory"`
	}{
		alias: (*alias)(sb),
	}
	err := json.Unmarshal(b, t)
	if err != nil {
		return err
	}

	if len(t.Accessory) > 0 {
		accessory, err := UnmarshalBlockElement(t.Accessory)
		if err != nil {
			return xerrors.Errorf("failed to unmarshal given element: %w", err)
		}
		sb.Accessory = accessory
	}

	return nil
}
