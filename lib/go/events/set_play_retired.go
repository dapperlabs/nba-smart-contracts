package events

import (
	"fmt"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/events/decoder"
)

const (
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
	cadenceValue, err := decoder.GetCadenceEvent(b)
	if err != nil {
		return nil, err
	}
	if cadenceValue.EventType.QualifiedIdentifier != EventPlayRetiredFromSet {
		return nil, fmt.Errorf("unexpected event type: %s", cadenceValue.EventType.QualifiedIdentifier)
	}
	eventMap, err := decoder.ConvertEvent(cadenceValue)
	event := setPlayRetiredEvent(eventMap)
	return event, nil
}
