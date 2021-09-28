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

func (evt playCreatedEvent) validate() error {
	if evt.EventType.QualifiedIdentifier != EventPlayCreated{
		return fmt.Errorf("error validating event: event is not a valid play created event, expected type %s, got %s",
			EventPlayCreated, evt.EventType.QualifiedIdentifier)
	}
	return nil
}

func DecodePlayCreatedEvent(b []byte) (PlayCreatedEvent, error) {
	value, err := jsoncdc.Decode(b)
	if err != nil {
		return nil, err
	}
	event := playCreatedEvent(value.(cadence.Event))
	if err := event.validate(); err != nil{
		return nil, fmt.Errorf("error decoding event: %w", err)
	}
	return event, nil
}
