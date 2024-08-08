package events

import (
	"fmt"
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

func (evt momentDestroyedEvent) validate() error {
	if evt["eventType"].(string) != EventMomentDestroyed {
		return fmt.Errorf("error validating event: event is not a valid moment destroyed event, expected type %s, got %s",
			EventMomentDestroyed, evt["eventType"].(string))
	}
	return nil
}

func DecodeMomentDestroyedEvent(b []byte) (MomentDestroyedEvent, error) {
	eventMap, err := decoder.DecodeToEventMap(b)
	if err != nil {
		return nil, err
	}
	event := momentDestroyedEvent(eventMap)
	return event, nil
}
