package events

import (
	"fmt"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/events/decoder"
)

const (
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
	cadenceValue, err := decoder.GetCadenceEvent(b)
	if err != nil {
		return nil, err
	}
	if cadenceValue.EventType.QualifiedIdentifier != EventSetCreated {
		return nil, fmt.Errorf("unexpected event type: %s", cadenceValue.EventType.QualifiedIdentifier)
	}
	eventMap, err := decoder.ConvertEvent(cadenceValue)
	event := setCreatedEvent(eventMap)
	return event, nil
}
