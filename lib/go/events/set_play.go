package events

import (
	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
)

var (
	EventPlayAddedToSet string = "TopShot.PlayAddedToSet"
)

type PlayAddedToSetEvent interface {
	SetID() uint32
	PlayID() uint32
}

type playAddedToSetEvent cadence.Event

func (p playAddedToSetEvent) SetID() uint32 {
	return uint32(p.Fields[0].(cadence.UInt32))
}

func (p playAddedToSetEvent) PlayID() uint32 {
	return uint32(p.Fields[1].(cadence.UInt32))
}

var _ PlayAddedToSetEvent = (*playAddedToSetEvent)(nil)

func DecodePlayAddedToSetEvent(b []byte)(PlayAddedToSetEvent, error) {
	value, err := jsoncdc.Decode(b)
	if err != nil {
		return nil, err
	}
	return playAddedToSetEvent(value.(cadence.Event)), nil

}
