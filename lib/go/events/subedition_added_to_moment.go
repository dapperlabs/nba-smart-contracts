package events

import (
	"fmt"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
)

var (
	EventSubeditionAddedToMoment = "TopShot.SubeditionAddedToMoment"
)

type SubeditionAddedToMomentEvent interface {
	MomentID() uint64
	SubeditionID() uint32
}

type subeditionAddedToMomentEvent cadence.Event

func (evt subeditionAddedToMomentEvent) MomentID() uint64 {
	return evt.Fields[0].(cadence.UInt64).ToGoValue().(uint64)
}

func (evt subeditionAddedToMomentEvent) SubeditionID() uint32 {
	return evt.Fields[1].(cadence.UInt32).ToGoValue().(uint32)
}

func (evt subeditionAddedToMomentEvent) validate() error {
	if evt.EventType.QualifiedIdentifier != EventSubeditionAddedToMoment {
		return fmt.Errorf("error validating event: event is not a valid subedition added to moment event, expected type %s, got %s",
			EventSubeditionAddedToMoment, evt.EventType.QualifiedIdentifier)
	}
	return nil
}

func DecodeSubeditionAddedToMomentEvent(b []byte) (SubeditionAddedToMomentEvent, error) {
	value, err := jsoncdc.Decode(nil, b)
	if err != nil {
		return nil, err
	}
	event := subeditionAddedToMomentEvent(value.(cadence.Event))
	if err := event.validate(); err != nil {
		return nil, fmt.Errorf("error decoding event: %w", err)
	}
	return event, nil
}
