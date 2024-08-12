package events

import (
	"fmt"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/events/decoder"
)

const (
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
	cadenceValue, err := decoder.GetCadenceEvent(b)
	if err != nil {
		return nil, err
	}
	if cadenceValue.EventType.QualifiedIdentifier != MomentUnlocked {
		return nil, fmt.Errorf("unexpected event type: %s", cadenceValue.EventType.QualifiedIdentifier)
	}
	eventMap, err := decoder.ConvertEvent(cadenceValue)

	event := momentUnlockedEvent(eventMap)
	return event, nil
}
