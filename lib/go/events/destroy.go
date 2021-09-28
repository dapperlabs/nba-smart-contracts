package events

import (
	"fmt"
	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
)

var (
	EventMomentDestroyed = "TopShot.MomentDestroyed"
)

type MomentDestroyedEvent interface {
	Id() uint64
}

type momentDestroyedEvent cadence.Event

func (evt momentDestroyedEvent) Id() uint64 {
	return uint64(evt.Fields[0].(cadence.UInt64))
}

func (evt momentDestroyedEvent) validate() error {
	if evt.EventType.QualifiedIdentifier != EventMomentDestroyed{
		return fmt.Errorf("error validating event: event is not a valid moment destroyed event, expected type %s, got %s",
			EventMomentDestroyed, evt.EventType.QualifiedIdentifier)
	}
	return nil
}

func DecodeMomentDestroyedEvent(b []byte) (MomentDestroyedEvent, error) {
	value, err := jsoncdc.Decode(b)
	if err != nil {
		return nil, err
	}
	event := momentDestroyedEvent(value.(cadence.Event))
	if err := event.validate(); err != nil{
		return nil, fmt.Errorf("error decoding event: %w", err)
	}
	return event, nil
}
