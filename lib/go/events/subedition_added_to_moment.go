package events

import (
	"fmt"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/events/decoder"
)

var (
	EventSubeditionAddedToMoment = "TopShot.SubeditionAddedToMoment"
)

type SubeditionAddedToMomentEvent interface {
	MomentID() uint64
	SubeditionID() uint32
}

type subeditionAddedToMomentEvent map[string]any

func (evt subeditionAddedToMomentEvent) MomentID() uint64 {
	return evt["momentID"].(uint64)
}

func (evt subeditionAddedToMomentEvent) SubeditionID() uint32 {
	return evt["subeditionID"].(uint32)
}

func (evt subeditionAddedToMomentEvent) validate() error {
	if evt["eventType"].(string) != EventSubeditionAddedToMoment {
		return fmt.Errorf("error validating event: event is not a valid subedition added to moment event, expected type %s, got %s",
			EventSubeditionAddedToMoment, evt["eventType"].(string))
	}
	return nil
}

func DecodeSubeditionAddedToMomentEvent(b []byte) (SubeditionAddedToMomentEvent, error) {
	eventMap, err := decoder.DecodeToEventMap(b)
	if err != nil {
		return nil, err
	}
	event := subeditionAddedToMomentEvent(eventMap)
	if err := event.validate(); err != nil {
		return nil, fmt.Errorf("error decoding event: %w", err)
	}
	return event, nil
}
