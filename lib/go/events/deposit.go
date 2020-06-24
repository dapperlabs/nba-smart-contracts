package events

import (
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

func DecodeDepositEvent(b []byte) (DepositEvent, error) {
	value, err := jsoncdc.Decode(b)
	if err != nil {
		return nil, err
	}
	return depositEvent(value.(cadence.Event)), nil
}
