package event

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

type FileID string

func (id FileID) String() string {
	return string(id)
}

type TeamID string

func (id TeamID) String() string {
	return string(id)
}

type SubTeamID string

func (id SubTeamID) String() string {
	return string(id)
}

type CommentID string

func (id CommentID) String() string {
	return string(id)
}
