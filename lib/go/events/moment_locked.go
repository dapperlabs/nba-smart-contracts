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
