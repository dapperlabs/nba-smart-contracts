package events

import (
	"fmt"
	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
)

var (
	EventPlayRetiredFromSet = "TopShot.PlayRetiredFromSet"
)

type SetPlayRetiredEvent interface {
	SetID() uint32
	PlayID() uint32
	NumMoments() uint32
}

type setPlayRetiredEvent cadence.Event

func (evt setPlayRetiredEvent) SetID() uint32 {
	return uint32(evt.Fields[0].(cadence.UInt32))
}

func (evt setPlayRetiredEvent) PlayID() uint32 {
	return uint32(evt.Fields[1].(cadence.UInt32))
}

func (evt setPlayRetiredEvent) NumMoments() uint32 {
	return uint32(evt.Fields[2].(cadence.UInt32))
}

func (evt setPlayRetiredEvent) isValidEvent() bool {
	return evt.EventType.QualifiedIdentifier == EventPlayRetiredFromSet
}

var _ SetPlayRetiredEvent = (*setPlayRetiredEvent)(nil)

func DecodeSetPlayRetiredEvent(b []byte) (SetPlayRetiredEvent, error) {
	value, err := jsoncdc.Decode(b)
	if err != nil {
		return nil, err
	}
	event := setPlayRetiredEvent(value.(cadence.Event))
	if !event.isValidEvent(){
		return nil, fmt.Errorf("error decoding event: event is not a valid play retired from set event")
	}
	return event, nil
}
