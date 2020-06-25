package events

import (
	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
)

var (
	EventPlayCreated string = "TopShot.PlayCreated"
)

type PlayCreatedEvent interface {
	Id() uint32
	MetaData() map[interface{}]interface{}
}

type playCreatedEvent cadence.Event


func (evt playCreatedEvent) Id() uint32 {
	return evt.Fields[0].(cadence.UInt32).ToGoValue().(uint32)
}
func (evt playCreatedEvent) MetaData() map[interface{}]interface{} {
	return evt.Fields[1].(cadence.Dictionary).ToGoValue().(map[interface{}]interface{})
}

func DecodePlayCreatedEvent(b []byte) (PlayCreatedEvent, error) {
	value, err := jsoncdc.Decode(b)
	if err != nil {
		return nil, err
	}
	return playCreatedEvent(value.(cadence.Event)), nil
}
