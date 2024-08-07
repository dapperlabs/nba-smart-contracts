package events

import (
	"fmt"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/events/decoder"
)

var (
	EventPlayRetiredFromSet = "TopShot.PlayRetiredFromSet"
)

type SetPlayRetiredEvent interface {
	SetID() uint32
	PlayID() uint32
	NumMoments() uint32
}

type setPlayRetiredEvent map[string]any

func (evt setPlayRetiredEvent) SetID() uint32 {
	return evt["setID"].(uint32)
}

func (evt setPlayRetiredEvent) PlayID() uint32 {
	return evt["playID"].(uint32)
}

func (evt setPlayRetiredEvent) NumMoments() uint32 {
	return evt["numMoments"].(uint32)
}

func (evt setPlayRetiredEvent) validate() error {
	if evt["eventType"].(string) != EventPlayRetiredFromSet {
		return fmt.Errorf("error validating event: event is not a valid play retired from set event, expected type %s, got %s",
			EventPlayRetiredFromSet, evt["eventType"].(string))
	}
	return nil
}

var _ SetPlayRetiredEvent = (*setPlayRetiredEvent)(nil)

func DecodeSetPlayRetiredEvent(b []byte) (SetPlayRetiredEvent, error) {
	eventMap, err := decoder.DecodeToEventMap(b)
	if err != nil {
		return nil, err
	}
	event := setPlayRetiredEvent(eventMap)
	if err := event.validate(); err != nil {
		return nil, fmt.Errorf("error decoding event: %w", err)
	}
	return event, nil
}
