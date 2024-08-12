package events

import (
	"fmt"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/events/decoder"
)

const (
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

func DecodeDepositEvent(b []byte) (DepositEvent, error) {
	cadenceValue, err := decoder.GetCadenceEvent(b)
	if err != nil {
		return nil, err
	}
	if cadenceValue.EventType.QualifiedIdentifier != EventDeposit {
		return nil, fmt.Errorf("unexpected event type: %s", cadenceValue.EventType.QualifiedIdentifier)
	}
	eventMap, err := decoder.ConvertEvent(cadenceValue)
	event := depositEvent(eventMap)
	return event, nil
}
