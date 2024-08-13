package events

import (
	"fmt"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/events/decoder"
)

var (
	MomentLocked = "TopShotLocking.MomentLocked"
)

type MomentLockedEvent interface {
	FlowID() uint64
	Duration() uint64
	ExpiryTimestamp() uint64
}

type momentLockedEvent map[string]any

func (evt momentLockedEvent) FlowID() uint64 {
	return evt["id"].(uint64)
}

func (evt momentLockedEvent) Duration() uint64 {
	return evt["duration"].(uint64)
}

func (evt momentLockedEvent) ExpiryTimestamp() uint64 {
	return evt["expiryTimestamp"].(uint64)
}

var _ MomentLockedEvent = (*momentLockedEvent)(nil)

func DecodeMomentLockedEvent(b []byte) (MomentLockedEvent, error) {
	cadenceValue, err := decoder.GetCadenceEvent(b)
	if err != nil {
		return nil, err
	}
	if cadenceValue.EventType.QualifiedIdentifier != MomentLocked {
		return nil, fmt.Errorf("unexpected event type: %s", cadenceValue.EventType.QualifiedIdentifier)
	}
	eventMap, err := decoder.ConvertEvent(cadenceValue)

	event := momentLockedEvent(eventMap)

	return event, nil
}
