package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nimona/go-nimona/repository"
)

type Pair struct {
	Key   []byte  `json:"k"`
	Value *[]byte `json:"v"`
}

func (kv *Pair) GetGUID() []byte {
	return kv.Key
}

func (kv *Pair) Marshal() ([]byte, error) {
	return json.Marshal(kv)
}

func (kv *Pair) Unmarshal(v []byte) error {
	return json.Unmarshal(v, kv)
}

type PairSetEventPayload struct {
	Value *[]byte `json:"value"`
}

func (kv *Pair) ApplyEvent(genericEvent repository.Event) {
	eventJSON, _ := genericEvent.Marshal() // TODO(geoah) Check error

	var payloadRawJSON json.RawMessage
	event := &Event{
		Payload: &payloadRawJSON,
	}
	json.Unmarshal(eventJSON, event) // TODO(geoah) Check error

	fmt.Printf("\t> Applying event. event=%#v;\n", event)
	switch event.GetTopic() {
	case "set":
		payload := &PairSetEventPayload{}
		if err := json.Unmarshal(payloadRawJSON, payload); err != nil {
			log.Fatal(err)
			return
		}
		kv.Key = event.GetGUID()
		kv.Value = payload.Value
		if kv.Value != nil {
			fmt.Printf("\t\tApplied SET event. key=%#v; value=%#v;\n", string(kv.Key), string(*kv.Value))
		} else {
			fmt.Printf("\t\tApplied DEL event. key=%#v;\n", string(kv.Key))
		}
	}
}
