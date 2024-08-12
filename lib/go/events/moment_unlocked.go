package events

import (
	"github.com/dapperlabs/nba-smart-contracts/lib/go/events/decoder"
)

var (
	MomentUnlocked = "TopShotLocking.MomentUnlocked"
)

type MomentUnlockedEvent interface {
	FlowID() uint64
}

type momentUnlockedEvent map[string]any

func (evt momentUnlockedEvent) FlowID() uint64 {
	return evt["flowID"].(uint64)
}

var _ MomentUnlockedEvent = (*momentUnlockedEvent)(nil)

func DecodeMomentUnlockedEvent(b []byte) (MomentUnlockedEvent, error) {
	eventMap, err := decoder.DecodeToEventMap(b)
	if err != nil {
		return nil, err
	}

	event := momentUnlockedEvent(eventMap)
	return event, nil
}
