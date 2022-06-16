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

func (evt setCreatedEvent) validate() error {
	if evt.EventType.QualifiedIdentifier != EventSetCreated {
		return fmt.Errorf("error validating event: event is not a valid set created event, expected type %s, got %s",
			EventSetCreated, evt.EventType.QualifiedIdentifier)
	}
	return nil
}

var _ SetCreatedEvent = (*setCreatedEvent)(nil)

func DecodeSetCreatedEvent(b []byte) (SetCreatedEvent, error) {
	value, err := jsoncdc.Decode(nil, b)
	if err != nil {
		return nil, err
	}
	event := setCreatedEvent(value.(cadence.Event))
	if err := event.validate(); err != nil {
		return nil, fmt.Errorf("error decoding event: %w", err)
	}
	return event, nil
}
