package events

import (
	"fmt"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/events/decoder"
)

const (
	EventMomentDestroyed   = "TopShot.MomentDestroyed"
	EventMomentDestroyedV2 = "TopShot.NFT.ResourceDestroyed"
)

type MomentDestroyedEvent interface {
	Id() uint64
}

type momentDestroyedEvent map[string]any

func (evt momentDestroyedEvent) Id() uint64 {
	return evt["id"].(uint64)
}

func DecodeMomentDestroyedEvent(b []byte) (MomentDestroyedEvent, error) {
	cadenceValue, err := decoder.GetCadenceEvent(b)
	if err != nil {
		return nil, err
	}
	if cadenceValue.EventType.QualifiedIdentifier != EventMomentDestroyed && cadenceValue.EventType.QualifiedIdentifier != EventMomentDestroyedV2 {
		return nil, fmt.Errorf("unexpected event type: %s", cadenceValue.EventType.QualifiedIdentifier)
	}
	eventMap, err := decoder.ConvertEvent(cadenceValue)
	event := momentDestroyedEvent(eventMap)
	return event, nil
}
