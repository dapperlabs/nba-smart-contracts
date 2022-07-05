package events

import (
	"fmt"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow-go-sdk"
)

var (
	// This variable specifies that there is a Withdraw Event on a TopShot Contract located at address 0x04
	EventWithdraw = "TopShot.Withdraw"
)

type WithdrawEvent interface {
	Id() uint64
	Owner() string // deprecated: use From()
	From() string
}

type withdrawEvent cadence.Event

var _ WithdrawEvent = (*withdrawEvent)(nil)

func (evt withdrawEvent) Id() uint64 {
	return uint64(evt.Fields[0].(cadence.UInt64))
}

func (evt withdrawEvent) From() string {
	optionalAddress := (evt.Fields[1]).(cadence.Optional)
	if cadenceAddress, ok := optionalAddress.Value.(cadence.Address); ok {
		return flow.BytesToAddress(cadenceAddress.Bytes()).Hex()
	}
	return "undefined"
}

func (evt withdrawEvent) Owner() string {
	return evt.From()
}

func (evt withdrawEvent) validate() error {
	if evt.EventType.QualifiedIdentifier != EventWithdraw {
		return fmt.Errorf("error validating event: event is not a valid withdraw event, expected type %s, got %s",
			EventWithdraw, evt.EventType.QualifiedIdentifier)
	}
	return nil
}

func DecodeWithdrawEvent(b []byte) (WithdrawEvent, error) {
	value, err := jsoncdc.Decode(nil, b)
	if err != nil {
		return nil, err
	}
	event := withdrawEvent(value.(cadence.Event))
	if err := event.validate(); err != nil {
		return nil, fmt.Errorf("error decoding event: %w", err)
	}
	return event, nil
}
