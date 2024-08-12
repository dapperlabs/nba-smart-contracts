package events

import (
	"github.com/dapperlabs/nba-smart-contracts/lib/go/events/decoder"
	"strings"
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

type revealedEvent map[string]any

var _ RevealedEvent = (*revealedEvent)(nil)

func (evt revealedEvent) Id() uint64 {
	return evt["id"].(uint64)
}

func (evt revealedEvent) Salt() string {
	return evt["salt"].(string)
}

func (evt revealedEvent) NFTs() string {
	return evt["nfts"].(string)
}

func parseNFTs(nft string) []string {
	return strings.Split(nft, ",")
}

func DecodeRevealedEvent(b []byte) (RevealedEvent, error) {
	eventMap, err := decoder.DecodeToEventMap(b)
	if err != nil {
		return nil, err
	}
	event := revealedEvent(eventMap)
	return event, nil
}
