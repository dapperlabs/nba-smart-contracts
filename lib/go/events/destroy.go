package events

import (
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

func (a momentDestroyedEvent) Id() uint64 {
	return uint64(a.Fields[0].(cadence.UInt64))
}

func DecodeMomentDestroyedEvent(b []byte) (MomentDestroyedEvent, error) {
	value, err := jsoncdc.Decode(b)
	if err != nil {
		return nil, err
	}
	return momentDestroyedEvent(value.(cadence.Event)), nil
}
