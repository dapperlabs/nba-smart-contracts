package events

import (
	"fmt"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
)

var (
	EventSubeditionCreated = "TopShot.SubeditionCreated"
)

type SubeditionCreatedEvent interface {
	Id() uint32
	Name() string
	MetaData() map[interface{}]interface{}
}

type subeditionCreatedEvent cadence.Event

func (evt subeditionCreatedEvent) Id() uint32 {
	return evt.Fields[0].(cadence.UInt32).ToGoValue().(uint32)
}

func (evt subeditionCreatedEvent) Name() string {
	return evt.Fields[1].(cadence.String).ToGoValue().(string)
}

func (evt subeditionCreatedEvent) MetaData() map[interface{}]interface{} {
	return evt.Fields[2].(cadence.Dictionary).ToGoValue().(map[interface{}]interface{})
}

func (evt subeditionCreatedEvent) validate() error {
	if evt.EventType.QualifiedIdentifier != EventSubeditionCreated {
		return fmt.Errorf("error validating event: event is not a valid subedition created event, expected type %s, got %s",
			EventSubeditionCreated, evt.EventType.QualifiedIdentifier)
	}
	return nil
}

func DecodeSubeditionCreatedEvent(b []byte) (SubeditionCreatedEvent, error) {
	value, err := jsoncdc.Decode(nil, b)
	if err != nil {
		return nil, err
	}
	event := subeditionCreatedEvent(value.(cadence.Event))
	if err := event.validate(); err != nil {
		return nil, fmt.Errorf("error decoding event: %w", err)
	}
	return event, nil
}
