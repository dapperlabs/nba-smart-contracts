package events

import (
	"fmt"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/events/decoder"
)

const (
	EventWithdraw = "TopShot.Withdraw"
)

type WithdrawEvent interface {
	Id() uint64
	From() string
	Owner() string
}

type withdrawEvent map[string]any

var _ WithdrawEvent = (*withdrawEvent)(nil)

func (evt withdrawEvent) Id() uint64 {
	return evt["id"].(uint64)
}

func (evt withdrawEvent) From() string {
	optionalAddress := evt["from"]
	if optionalAddress == nil {
		return "undefined"
	}
	return optionalAddress.(string)
}

func (evt withdrawEvent) Owner() string {
	return evt.From()
}

func DecodeWithdrawEvent(b []byte) (WithdrawEvent, error) {
	cadenceValue, err := decoder.GetCadenceEvent(b)
	if err != nil {
		return nil, err
	}
	if cadenceValue.EventType.QualifiedIdentifier != EventWithdraw {
		return nil, fmt.Errorf("unexpected event type: %s", cadenceValue.EventType.QualifiedIdentifier)
	}
	eventMap, err := decoder.ConvertEvent(cadenceValue)
	event := withdrawEvent(eventMap)
	return event, nil
}
