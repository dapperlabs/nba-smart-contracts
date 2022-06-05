package events

import (
	"fmt"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow-go-sdk"
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

type depositEvent cadence.Event

var _ DepositEvent = (*depositEvent)(nil)

func (evt depositEvent) Id() uint64 {
	return uint64(evt.Fields[0].(cadence.UInt64))
}

func (evt depositEvent) To() string {
	optionalAddress := (evt.Fields[1]).(cadence.Optional)
	if cadenceAddress, ok := optionalAddress.Value.(cadence.Address); ok {
		return flow.BytesToAddress(cadenceAddress.Bytes()).Hex()
	}
	return "undefined"
}

func (evt depositEvent) Owner() string {
	return evt.To()
}

func (evt depositEvent) validate() error {
	if evt.EventType.QualifiedIdentifier != EventDeposit {
		return fmt.Errorf("error validating event: event is not a valid moment destroyed event, expected type %s, got %s",
			EventDeposit, evt.EventType.QualifiedIdentifier)
	}
	return nil
}

func DecodeDepositEvent(b []byte) (DepositEvent, error) {
	value, err := jsoncdc.Decode(nil, b)
	if err != nil {
		return nil, err
	}
	event := depositEvent(value.(cadence.Event))
	if err := event.validate(); err != nil {
		return nil, fmt.Errorf("error decoding event: %w", err)
	}
	return event, nil
}
