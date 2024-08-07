package events

import (
	"fmt"
	"github.com/dapperlabs/nba-smart-contracts/lib/go/events/decoder"
)

var (
	EventSubeditionCreated = "TopShot.SubeditionCreated"
)

type SubeditionCreatedEvent interface {
	SubeditionId() uint32
	Name() string
	MetaData() map[string]interface{}
}

type subeditionCreatedEvent map[string]any

func (evt subeditionCreatedEvent) SubeditionId() uint32 {
	return evt["subeditionId"].(uint32)
}

func (evt subeditionCreatedEvent) Name() string {
	return evt["name"].(string)
}

func (evt subeditionCreatedEvent) MetaData() map[string]interface{} {
	metadata := evt["metadata"].(map[interface{}]interface{})
	result := make(map[string]interface{})
	for k, v := range metadata {
		result[k.(string)] = v
	}
	return result
}

func (evt subeditionCreatedEvent) validate() error {
	if evt["eventType"].(string) != EventSubeditionCreated {
		return fmt.Errorf("error validating event: event is not a valid subedition created event, expected type %s, got %s",
			EventSubeditionCreated, evt["eventType"].(string))
	}
	return nil
}

func DecodeSubeditionCreatedEvent(b []byte) (SubeditionCreatedEvent, error) {
	eventMap, err := decoder.DecodeToEventMap(b)
	if err != nil {
		return nil, err
	}
	event := subeditionCreatedEvent(eventMap)
	if err := event.validate(); err != nil {
		return nil, fmt.Errorf("error decoding event: %w", err)
	}
	return event, nil
}
