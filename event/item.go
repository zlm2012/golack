package event

type ItemType string

const (
	ItemTypeChannelMessage ItemType = "C"
	ItemTypeGroupMessage   ItemType = "G"
	ItemTypeFile           ItemType = "F"
	ItemTypeFileComment    ItemType = "Fc"
)
