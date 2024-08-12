package events

import (
	"fmt"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/events/decoder"
)

const (
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

var _ PlayAddedToSetEvent = (*playAddedToSetEvent)(nil)

func DecodePlayAddedToSetEvent(b []byte) (PlayAddedToSetEvent, error) {
	cadenceValue, err := decoder.GetCadenceEvent(b)
	if err != nil {
		return nil, err
	}
	if cadenceValue.EventType.QualifiedIdentifier != EventPlayAddedToSet {
		return nil, fmt.Errorf("unexpected event type: %s", cadenceValue.EventType.QualifiedIdentifier)
	}
	eventMap, err := decoder.ConvertEvent(cadenceValue)
	event := playAddedToSetEvent(eventMap)
	return event, nil
}
