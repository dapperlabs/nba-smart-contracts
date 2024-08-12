package events

import (
	"github.com/dapperlabs/nba-smart-contracts/lib/go/events/decoder"
)

var (
	EventSetCreated = "TopShot.SetCreated"
)

type SetCreatedEvent interface {
	SetID() uint32
	Series() uint32
}

type setCreatedEvent map[string]any

func (evt setCreatedEvent) SetID() uint32 {
	return evt["setID"].(uint32)
}

func (evt setCreatedEvent) Series() uint32 {
	return evt["series"].(uint32)
}

var _ SetCreatedEvent = (*setCreatedEvent)(nil)

func DecodeSetCreatedEvent(b []byte) (SetCreatedEvent, error) {
	eventMap, err := decoder.DecodeToEventMap(b)
	if err != nil {
		return nil, err
	}
	event := setCreatedEvent(eventMap)
	return event, nil
}
