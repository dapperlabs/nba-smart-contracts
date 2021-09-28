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

func (evt momentDestroyedEvent) isValidEvent() bool {
	return evt.EventType.QualifiedIdentifier == EventMomentDestroyed
}

func DecodeMomentDestroyedEvent(b []byte) (MomentDestroyedEvent, error) {
	value, err := jsoncdc.Decode(b)
	if err != nil {
		return nil, err
	}
	event := momentDestroyedEvent(value.(cadence.Event))
	if !event.isValidEvent(){
		return nil, fmt.Errorf("error decoding event: event is not a valid deposit event")
	}
	return event, nil
}
