package events

import (
	"github.com/dapperlabs/nba-smart-contracts/lib/go/events/decoder"
)

var (
	MomentLocked = "TopShotLocking.MomentLocked"
)

type MomentLockedEvent interface {
	FlowID() uint64
	Duration() float64
	ExpiryTimestamp() float64
}

type momentLockedEvent map[string]any

func (evt momentLockedEvent) FlowID() uint64 {
	return evt["flowID"].(uint64)
}

func (evt momentLockedEvent) Duration() float64 {
	return evt["duration"].(float64)
}

func (evt momentLockedEvent) ExpiryTimestamp() float64 {
	return evt["expiryTimestamp"].(float64)
}

var _ MomentLockedEvent = (*momentLockedEvent)(nil)

func DecodeMomentLockedEvent(b []byte) (MomentLockedEvent, error) {
	eventMap, err := decoder.DecodeToEventMap(b)
	if err != nil {
		return nil, err
	}

	event := momentLockedEvent(eventMap)

	return event, nil
}
