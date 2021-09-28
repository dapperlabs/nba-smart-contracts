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

func (evt setLockedEvent) validate() error {
	if evt.EventType.QualifiedIdentifier != EventSetLocked{
		return fmt.Errorf("error validating event: event is not a valid set locked event, expected type %s, got %s",
			EventSetLocked, evt.EventType.QualifiedIdentifier)
	}
	return nil
}

func DecodeSetLockedEvent(b []byte) (SetLockedEvent, error) {
	value, err := jsoncdc.Decode(b)
	if err != nil {
		return nil, err
	}
	event := setLockedEvent(value.(cadence.Event))
	if err := event.validate(); err != nil{
		return nil, fmt.Errorf("error decoding event: %w", err)
	}
	return event, nil
}
