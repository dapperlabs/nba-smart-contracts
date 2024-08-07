package events

import (
	"fmt"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/events/decoder"
)

var (
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
	return evt["metaData"].(map[interface{}]interface{})
}

func (evt playCreatedEvent) validate() error {
	if evt["eventType"].(string) != EventPlayCreated {
		return fmt.Errorf("error validating event: event is not a valid play created event, expected type %s, got %s",
			EventPlayCreated, evt["eventType"].(string))
	}
	return nil
}

func DecodePlayCreatedEvent(b []byte) (PlayCreatedEvent, error) {
	eventMap, err := decoder.DecodeToEventMap(b)
	if err != nil {
		return nil, err
	}
	event := playCreatedEvent(eventMap)
	if err := event.validate(); err != nil {
		return nil, fmt.Errorf("error decoding event: %w", err)
	}
	return event, nil
}
