package events

import (
	"fmt"
	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
)

var (
	EventPlayAddedToSet = "TopShot.PlayAddedToSet"
)

type PlayAddedToSetEvent interface {
	SetID() uint32
	PlayID() uint32
}

type playAddedToSetEvent cadence.Event

func (evt playAddedToSetEvent) SetID() uint32 {
	return uint32(evt.Fields[0].(cadence.UInt32))
}

func (evt playAddedToSetEvent) PlayID() uint32 {
	return uint32(evt.Fields[1].(cadence.UInt32))
}

func (evt playAddedToSetEvent) validate() error {
	if evt.EventType.QualifiedIdentifier != EventPlayAddedToSet{
		return fmt.Errorf("error validating event: event is not a valid play added to set event, expected type %s, got %s",
			EventPlayAddedToSet, evt.EventType.QualifiedIdentifier)
	}
	return nil
}

var _ PlayAddedToSetEvent = (*playAddedToSetEvent)(nil)

func DecodePlayAddedToSetEvent(b []byte)(PlayAddedToSetEvent, error) {
	value, err := jsoncdc.Decode(b)
	if err != nil {
		return nil, err
	}
	event := playAddedToSetEvent(value.(cadence.Event))
	if err := event.validate(); err != nil{
		return nil, fmt.Errorf("error decoding event: %w", err)
	}
	return event, nil

}
