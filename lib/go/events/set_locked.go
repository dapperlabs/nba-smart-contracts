package events

import (
	"fmt"
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

func (evt setLockedEvent) validate() error {
	if evt["eventType"].(string) != EventSetLocked {
		return fmt.Errorf("error validating event: event is not a valid set locked event, expected type %s, got %s",
			EventSetLocked, evt["eventType"].(string))
	}
	return nil
}

func DecodeSetLockedEvent(b []byte) (SetLockedEvent, error) {
	eventMap, err := decoder.DecodeToEventMap(b)
	if err != nil {
		return nil, err
	}
	event := setLockedEvent(eventMap)
	if err := event.validate(); err != nil {
		return nil, fmt.Errorf("error decoding event: %w", err)
	}
	return event, nil
}
