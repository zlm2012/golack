package event

type ConversationType string

const (
	ConversationTypeIM           ConversationType = "im"
	ConversationTypeMultiPartyIM ConversationType = "mpim"
	ConversationTypePrivate      ConversationType = "private"
	ConversationTypePublic       ConversationType = "public"
)
