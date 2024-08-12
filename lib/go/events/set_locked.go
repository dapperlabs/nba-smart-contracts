package events

import (
	"fmt"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/events/decoder"
)

const (
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
	cadenceValue, err := decoder.GetCadenceEvent(b)
	if err != nil {
		return nil, err
	}
	if cadenceValue.EventType.QualifiedIdentifier != EventSetLocked {
		return nil, fmt.Errorf("unexpected event type: %s", cadenceValue.EventType.QualifiedIdentifier)
	}
	eventMap, err := decoder.ConvertEvent(cadenceValue)
	event := setLockedEvent(eventMap)
	return event, nil
}
