package events

import (
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
	return evt["metadata"].(map[interface{}]interface{})
}

func DecodePlayCreatedEvent(b []byte) (PlayCreatedEvent, error) {
	eventMap, err := decoder.DecodeToEventMap(b)
	if err != nil {
		return nil, err
	}
	event := playCreatedEvent(eventMap)
	return event, nil
}
