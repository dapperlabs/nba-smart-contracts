package events

import (
	"fmt"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/events/decoder"
)

const (
	EventPlayCreated = "TopShot.PlayCreated"
)

type PlayCreatedEvent interface {
	Id() uint32
	MetaData() map[interface{}]interface{}
}

type playCreatedEvent map[string]any

func (evt playCreatedEvent) Id() uint32 {
	return evt["id"].(uint32)
}

func (evt playCreatedEvent) MetaData() map[interface{}]interface{} {
	return evt["metadata"].(map[interface{}]interface{})
}

func DecodePlayCreatedEvent(b []byte) (PlayCreatedEvent, error) {
	cadenceValue, err := decoder.GetCadenceEvent(b)
	if err != nil {
		return nil, err
	}
	if cadenceValue.EventType.QualifiedIdentifier != EventPlayCreated {
		return nil, fmt.Errorf("unexpected event type: %s", cadenceValue.EventType.QualifiedIdentifier)
	}
	eventMap, err := decoder.ConvertEvent(cadenceValue)
	event := playCreatedEvent(eventMap)
	return event, nil
}
