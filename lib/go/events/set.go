package events

import (
	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
)

var (
	EventSetCreated string = "TopShot.SetCreated"
)

type SetCreatedEvent interface {
	SetID() uint32
	Series() uint32
}

type setCreatedEvent cadence.Event

func (s setCreatedEvent) SetID() uint32 {
	return uint32(s.Fields[0].(cadence.UInt32))
}

func (s setCreatedEvent) Series() uint32 {
	return uint32(s.Fields[1].(cadence.UInt32))
}

var _ SetCreatedEvent = (*setCreatedEvent)(nil)

func DecodeSetCreatedEvent(b []byte) (SetCreatedEvent, error) {
	value, err := jsoncdc.Decode(b)
	if err != nil {
		return nil, err
	}
	return setCreatedEvent(value.(cadence.Event)), nil
}
