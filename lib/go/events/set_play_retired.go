package events

import (
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

var _ SetPlayRetiredEvent = (*setPlayRetiredEvent)(nil)

func DecodeSetPlayRetiredEvent(b []byte) (SetPlayRetiredEvent, error) {
	eventMap, err := decoder.DecodeToEventMap(b)
	if err != nil {
		return nil, err
	}
	event := setPlayRetiredEvent(eventMap)
	return event, nil
}
