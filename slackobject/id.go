package slackobject

type AppID string

func (id AppID) String() string {
	return string(id)
}

type BotID string

func (id BotID) String() string {
	return string(id)
}

type ChannelID string

func (id ChannelID) String() string {
	return string(id)
}

type UserID string

func (id UserID) String() string {
	return string(id)
}
