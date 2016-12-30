package rtmapi

import (
	"sync"
	"testing"
)

func TestNewOutgoingEventID(t *testing.T) {
	eventID := NewOutgoingEventID()

	if eventID.id != 0 {
		t.Errorf("id value is not starting from 0. id was %d", eventID.id)
		return
	}
}

func TestOutgoingEventID_Next(t *testing.T) {
	eventID := OutgoingEventID{
		id:    0,
		mutex: &sync.Mutex{},
	}

	nextID := eventID.Next()
	if nextID != 1 {
		t.Errorf("id 1 must be given on first Next() call. id was %d.", nextID)
		return
	}
}
