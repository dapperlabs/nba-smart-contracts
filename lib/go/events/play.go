package events

import (
	"fmt"
	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
)

var (
	EventPlayCreated = "TopShot.PlayCreated"
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

func (evt playCreatedEvent) isValidEvent() bool {
	return evt.EventType.QualifiedIdentifier == EventPlayCreated
}

func DecodePlayCreatedEvent(b []byte) (PlayCreatedEvent, error) {
	value, err := jsoncdc.Decode(b)
	if err != nil {
		return nil, err
	}
	event := playCreatedEvent(value.(cadence.Event))
	if !event.isValidEvent(){
		return nil, fmt.Errorf("error decoding event: event is not a valid play created event")
	}
	return event, nil
}
