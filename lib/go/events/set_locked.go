package events

import (
	"fmt"
	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
)

var (
	EventSetLocked = "TopShot.SetLocked"
)

type SetLockedEvent interface {
	SetID() uint32
}

type setLockedEvent cadence.Event

var _ SetLockedEvent = (*setLockedEvent)(nil)

func (evt setLockedEvent) SetID() uint32 {
	return uint32(evt.Fields[0].(cadence.UInt32))
}

func (evt setLockedEvent) isValidEvent() bool {
	return evt.EventType.QualifiedIdentifier == EventSetLocked
}

func DecodeSetLockedEvent(b []byte) (SetLockedEvent, error) {
	value, err := jsoncdc.Decode(b)
	if err != nil {
		return nil, err
	}
	event := setLockedEvent(value.(cadence.Event))
	if !event.isValidEvent(){
		return nil, fmt.Errorf("error decoding event: event is not a valid set locked event")
	}
	return event, nil
}
