package event

type PresenceQuery struct {
	TypedEvent
	UserIDs []UserID `json:"ids"`
}

type PresenceSubscribe struct {
	TypedEvent
	UserIDs []UserID `json:"ids"`
}
