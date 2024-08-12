package events

import (
	"github.com/dapperlabs/nba-smart-contracts/lib/go/events/decoder"
)

var (
	EventSetLocked = "TopShot.SetLocked"
)

type SetLockedEvent interface {
	SetID() uint32
}

type setLockedEvent map[string]any

var _ SetLockedEvent = (*setLockedEvent)(nil)

func (evt setLockedEvent) SetID() uint32 {
	return evt["setID"].(uint32)
}

func DecodeSetLockedEvent(b []byte) (SetLockedEvent, error) {
	eventMap, err := decoder.DecodeToEventMap(b)
	if err != nil {
		return nil, err
	}
	event := setLockedEvent(eventMap)
	return event, nil
}
