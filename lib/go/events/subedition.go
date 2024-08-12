package events

import (
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
	return evt["subeditionID"].(uint32)
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

func DecodeSubeditionCreatedEvent(b []byte) (SubeditionCreatedEvent, error) {
	eventMap, err := decoder.DecodeToEventMap(b)
	if err != nil {
		return nil, err
	}
	event := subeditionCreatedEvent(eventMap)
	return event, nil
}
