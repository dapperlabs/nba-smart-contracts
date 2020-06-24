package events

import (
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

func (a momentMintedEvent) MomentId() uint64 {
	return uint64(a.Fields[0].(cadence.UInt64))
}

func (a momentMintedEvent) PlayId() uint32 {
	return uint32(a.Fields[1].(cadence.UInt32))
}

func (a momentMintedEvent) SetId() uint32 {
	return uint32(a.Fields[2].(cadence.UInt32))
}

func (a momentMintedEvent) SerialNumber() uint32 {
	return uint32(a.Fields[3].(cadence.UInt32))
}

var _ MomentMintedEvent = (*momentMintedEvent)(nil)

func DecodeMomentMintedEvent(b []byte) (MomentMintedEvent, error) {
	value, err := jsoncdc.Decode(b)
	if err != nil {
		return nil, err
	}
	return momentMintedEvent(value.(cadence.Event)), nil
}
