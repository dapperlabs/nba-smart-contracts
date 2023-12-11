package events

import (
	"fmt"
	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
)

var (
	MomentUnlocked = "TopShotLocking.MomentUnlocked"
)

type MomentUnlockedEvent interface {
	FlowID() uint64
}

type momentUnlockedEvent cadence.Event

func (evt momentUnlockedEvent) FlowID() uint64 {
	return uint64(evt.Fields[0].(cadence.UInt64))
}

func (evt momentUnlockedEvent) validate() error {
	if evt.EventType.QualifiedIdentifier != MomentUnlocked {
		return fmt.Errorf("error validating event: event is not a valid moment unlocked event, expected type %s, got %s",
			MomentUnlocked, evt.EventType.QualifiedIdentifier)
	}
	return nil
}

var _ MomentUnlockedEvent = (*momentUnlockedEvent)(nil)

func DecodeMomentUnlockedEvent(b []byte) (MomentUnlockedEvent, error) {
	value, err := jsoncdc.Decode(nil, b)
	if err != nil {
		return nil, err
	}

	event := momentUnlockedEvent(value.(cadence.Event))
	if err := event.validate(); err != nil {
		return nil, fmt.Errorf("error decoding event: %w", err)
	}

	return event, nil
}
