package events

import (
	"fmt"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/events/decoder"
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
	SubeditionId() uint32
}

type momentMintedEvent map[string]any

func (evt momentMintedEvent) MomentId() uint64 {
	return evt["momentId"].(uint64)
}

func (evt momentMintedEvent) PlayId() uint32 {
	return evt["playId"].(uint32)
}

func (evt momentMintedEvent) SetId() uint32 {
	return evt["setId"].(uint32)
}

func (evt momentMintedEvent) SerialNumber() uint32 {
	return evt["serialNumber"].(uint32)
}

func (evt momentMintedEvent) SubeditionId() uint32 {
	if val, ok := evt["subeditionId"]; ok {
		return val.(uint32)
	}
	return 0
}

func (evt momentMintedEvent) validate() error {
	if evt["eventType"].(string) != EventMomentMinted {
		return fmt.Errorf("error validating event: event is not a valid moment minted event, expected type %s, got %s",
			EventMomentMinted, evt["eventType"].(string))
	}
	return nil
}

var _ MomentMintedEvent = (*momentMintedEvent)(nil)

func DecodeMomentMintedEvent(b []byte) (MomentMintedEvent, error) {
	eventMap, err := decoder.DecodeToEventMap(b)
	if err != nil {
		return nil, err
	}
	event := momentMintedEvent(eventMap)
	if err := event.validate(); err != nil {
		return nil, fmt.Errorf("error decoding event: %w", err)
	}
	return event, nil
}