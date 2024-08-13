package events

import (
	"fmt"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/events/decoder"
)

const (
	// This variable specifies that there is a MomentMinted Event on a TopShot Contract located at address 0x04
	EventMomentMinted = "TopShot.MomentMinted"
)

type MomentMintedEvent interface {
	MomentId() uint64
	PlayId() uint32
	SetId() uint32
	SerialNumber() uint32
	SubeditionId() uint32
}

type momentMintedEvent map[string]any

func (evt momentMintedEvent) MomentId() uint64 {
	return evt["momentID"].(uint64)
}

func (evt momentMintedEvent) PlayId() uint32 {
	return evt["playID"].(uint32)
}

func (evt momentMintedEvent) SetId() uint32 {
	return evt["setID"].(uint32)
}

func (evt momentMintedEvent) SerialNumber() uint32 {
	return evt["serialNumber"].(uint32)
}

func (evt momentMintedEvent) SubeditionId() uint32 {
	if val, ok := evt["subeditionID"]; ok {
		return val.(uint32)
	}
	return 0
}

var _ MomentMintedEvent = (*momentMintedEvent)(nil)

func DecodeMomentMintedEvent(b []byte) (MomentMintedEvent, error) {
	cadenceValue, err := decoder.GetCadenceEvent(b)
	if err != nil {
		return nil, err
	}
	if cadenceValue.EventType.QualifiedIdentifier != EventMomentMinted {
		return nil, fmt.Errorf("unexpected event type: %s", cadenceValue.EventType.QualifiedIdentifier)
	}
	eventMap, err := decoder.ConvertEvent(cadenceValue)
	event := momentMintedEvent(eventMap)
	return event, nil
}
