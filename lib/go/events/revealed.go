package events

import (
	"fmt"
	"strings"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
)

const (
	// EventRevealed specifies that there is a Revealed Event on a PackNFT Contract located at the address
	EventRevealed = "PackNFT.Revealed"
)

type RevealedEvent interface {
	Id() uint64
	Salt() string
	NFTs() string
}

type revealedEvent cadence.Event

var _ RevealedEvent = (*revealedEvent)(nil)

func (evt revealedEvent) Id() uint64 {
	return uint64(evt.Fields[0].(cadence.UInt64))
}

func (evt revealedEvent) Salt() string {
	return string(evt.Fields[1].(cadence.String))
}

func (evt revealedEvent) NFTs() string {
	return string(evt.Fields[2].(cadence.String))
}

func (evt revealedEvent) validate() error {
	if evt.EventType.QualifiedIdentifier != EventRevealed {
		return fmt.Errorf("error validating event: event is not a valid revealed event, expected type %s, got %s",
			EventRevealed, evt.EventType.QualifiedIdentifier)
	}
	return nil
}

func parseNFTs(nft string) []string {
	return strings.Split(nft, ",")
}

func DecodeRevealedEvent(b []byte) (RevealedEvent, error) {
	value, err := jsoncdc.Decode(nil, b)
	if err != nil {
		return nil, err
	}
	event := revealedEvent(value.(cadence.Event))
	if err := event.validate(); err != nil {
		return nil, fmt.Errorf("error decoding event: %w", err)
	}
	return event, nil
}
