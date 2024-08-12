package events

import (
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

func DecodeSubeditionAddedToMomentEvent(b []byte) (SubeditionAddedToMomentEvent, error) {
	eventMap, err := decoder.DecodeToEventMap(b)
	if err != nil {
		return nil, err
	}
	event := subeditionAddedToMomentEvent(eventMap)
	return event, nil
}
