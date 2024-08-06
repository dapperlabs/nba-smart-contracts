package events

import (
	"fmt"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/events/decoder"
)

var (
	// This variable specifies that there is a Deposit Event on a TopShot Contract located at address 0x04
	EventDeposit = "TopShot.Deposit"
)

type DepositEvent interface {
	Id() uint64
	Owner() string // deprecated: use To()
	To() string
}

type depositEvent map[string]any

var _ DepositEvent = (*depositEvent)(nil)

func (evt depositEvent) Id() uint64 {
	return evt["id"].(uint64)
}

func (evt depositEvent) To() string {
	optionalAddress := evt["to"]
	if optionalAddress == nil {
		return "undefined"
	}
	address := optionalAddress.(string)
	return address
}

func (evt depositEvent) Owner() string {
	return evt.To()
}

func (evt depositEvent) validate() error {
	if evt["eventType"].(string) != EventDeposit {
		return fmt.Errorf("error validating event: event is not a valid moment destroyed event, expected type %s, got %s",
			EventDeposit, evt["eventType"].(string))
	}
	return nil
}

func DecodeDepositEvent(b []byte) (DepositEvent, error) {
	eventMap, err := decoder.DecodeToEventMap(b)
	if err != nil {
		return nil, err
	}
	event := depositEvent(eventMap)
	if err := event.validate(); err != nil {
		return nil, fmt.Errorf("error decoding event: %w", err)
	}
	return event, nil
}
