package events

import (
	"fmt"
	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
)

var (
	MomentLocked = "TopShotLocking.MomentLocked"
)

type MomentLockedEvent interface {
	FlowID() uint64
	Duration() uint64
	ExpiryTimestamp() uint64
}

type momentLockedEvent cadence.Event

func (evt momentLockedEvent) FlowID() uint64 {
	return uint64(evt.Fields[0].(cadence.UInt64))
}

func (evt momentLockedEvent) Duration() uint64 {
	return uint64(evt.Fields[1].(cadence.UInt64))
}

func (evt momentLockedEvent) ExpiryTimestamp() uint64 {
	return uint64(evt.Fields[2].(cadence.UInt64))
}

func (evt momentLockedEvent) validate() error {
	if evt.EventType.QualifiedIdentifier != MomentLocked {
		return fmt.Errorf("error validating event: event is not a valid moment locked event, expected type %s, got %s",
			MomentLocked, evt.EventType.QualifiedIdentifier)
	}
	return nil
}

var _ MomentLockedEvent = (*momentLockedEvent)(nil)

func DecodeMomentLockedEvent(b []byte) (MomentLockedEvent, error) {
	value, err := jsoncdc.Decode(nil, b)
	if err != nil {
		return nil, err
	}

	event := momentLockedEvent(value.(cadence.Event))
	if err := event.validate(); err != nil {
		return nil, fmt.Errorf("error decoding event: %w", err)
	}

	return event, nil
}
