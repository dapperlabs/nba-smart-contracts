package events

import (
	"fmt"
	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
)

var (
	EventSetCreated = "TopShot.SetCreated"
)

type SetCreatedEvent interface {
	SetID() uint32
	Series() uint32
}

type setCreatedEvent cadence.Event

func (evt setCreatedEvent) SetID() uint32 {
	return uint32(evt.Fields[0].(cadence.UInt32))
}

func (evt setCreatedEvent) Series() uint32 {
	return uint32(evt.Fields[1].(cadence.UInt32))
}

func (evt setCreatedEvent) isValidEvent() bool {
	return evt.EventType.QualifiedIdentifier == EventSetCreated
}

var _ SetCreatedEvent = (*setCreatedEvent)(nil)

func DecodeSetCreatedEvent(b []byte) (SetCreatedEvent, error) {
	value, err := jsoncdc.Decode(b)
	if err != nil {
		return nil, err
	}
	event := setCreatedEvent(value.(cadence.Event))
	if !event.isValidEvent(){
		return nil, fmt.Errorf("error decoding event: event is not a valid set created event")
	}
	return event, nil
}
