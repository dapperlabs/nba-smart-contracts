package events

import (
	"fmt"
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

func (evt momentUnlockedEvent) validate() error {
	if evt["eventType"].(string) != MomentUnlocked {
		return fmt.Errorf("error validating event: event is not a valid moment unlocked event, expected type %s, got %s",
			MomentUnlocked, evt["eventType"].(string))
	}
	return nil
}

var _ MomentUnlockedEvent = (*momentUnlockedEvent)(nil)

func DecodeMomentUnlockedEvent(b []byte) (MomentUnlockedEvent, error) {
	eventMap, err := decoder.DecodeToEventMap(b)
	if err != nil {
		return nil, err
	}

	event := momentUnlockedEvent(eventMap)
	if err := event.validate(); err != nil {
		return nil, fmt.Errorf("error decoding event: %w", err)
	}

	return event, nil
}
