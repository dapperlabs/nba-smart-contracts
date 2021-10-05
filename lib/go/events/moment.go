package events

import (
	"fmt"
	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
)

var (
	// This variable specifies that there is a MomentMinted Event on a TopShot Contract located at address 0x04
	EventMomentMinted = "TopShot.MomentMinted"
)

type MomentMintedEvent interface {
	MomentId() uint64
	PlayId() uint32
	SetId() uint32
	SerialNumber() uint32
}

type momentMintedEvent cadence.Event

func (evt momentMintedEvent) MomentId() uint64 {
	return uint64(evt.Fields[0].(cadence.UInt64))
}

func (evt momentMintedEvent) PlayId() uint32 {
	return uint32(evt.Fields[1].(cadence.UInt32))
}

func (evt momentMintedEvent) SetId() uint32 {
	return uint32(evt.Fields[2].(cadence.UInt32))
}

func (evt momentMintedEvent) SerialNumber() uint32 {
	return uint32(evt.Fields[3].(cadence.UInt32))
}

func (evt momentMintedEvent) validate() error {
	if evt.EventType.QualifiedIdentifier != EventMomentMinted{
		return fmt.Errorf("error validating event: event is not a valid moment minted event, expected type %s, got %s",
			EventMomentMinted, evt.EventType.QualifiedIdentifier)
	}
	return nil
}

var _ MomentMintedEvent = (*momentMintedEvent)(nil)

func DecodeMomentMintedEvent(b []byte) (MomentMintedEvent, error) {
	value, err := jsoncdc.Decode(b)
	if err != nil {
		return nil, err
	}
	event := momentMintedEvent(value.(cadence.Event))
	if err := event.validate(); err != nil{
		return nil, fmt.Errorf("error decoding event: %w", err)
	}
	return event, nil
}
