package events

import (
	"github.com/dapperlabs/nba-smart-contracts/lib/go/events/decoder"
)

var (
	EventMomentDestroyed = "TopShot.MomentDestroyed"
)

type MomentDestroyedEvent interface {
	Id() uint64
}

type momentDestroyedEvent map[string]any

func (evt momentDestroyedEvent) Id() uint64 {
	return evt["id"].(uint64)
}

func DecodeMomentDestroyedEvent(b []byte) (MomentDestroyedEvent, error) {
	eventMap, err := decoder.DecodeToEventMap(b)
	if err != nil {
		return nil, err
	}
	event := momentDestroyedEvent(eventMap)
	return event, nil
}
