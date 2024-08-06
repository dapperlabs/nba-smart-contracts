package events

import (
	"fmt"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/events/decoder"
)

var (
	EventPlayAddedToSet = "TopShot.PlayAddedToSet"
)

type PlayAddedToSetEvent interface {
	SetID() uint32
	PlayID() uint32
}

type playAddedToSetEvent map[string]any

func (evt playAddedToSetEvent) SetID() uint32 {
	return evt["setID"].(uint32)
}

func (evt playAddedToSetEvent) PlayID() uint32 {
	return evt["playID"].(uint32)
}

func (evt playAddedToSetEvent) validate() error {
	if evt["eventType"].(string) != EventPlayAddedToSet {
		return fmt.Errorf("error validating event: event is not a valid play added to set event, expected type %s, got %s",
			EventPlayAddedToSet, evt["eventType"].(string))
	}
	return nil
}

var _ PlayAddedToSetEvent = (*playAddedToSetEvent)(nil)

func DecodePlayAddedToSetEvent(b []byte) (PlayAddedToSetEvent, error) {
	eventMap, err := decoder.DecodeToEventMap(b)
	if err != nil {
		return nil, err
	}
	event := playAddedToSetEvent(eventMap)
	if err := event.validate(); err != nil {
		return nil, fmt.Errorf("error decoding event: %w", err)
	}
	return event, nil
}
