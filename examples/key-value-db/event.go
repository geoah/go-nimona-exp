package main

import "encoding/json"

type Event struct {
	Guid    []byte      `json:"guid"`
	Topic   string      `json:"topic"`
	Payload interface{} `json:"payload"`
}

func (e *Event) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

func (e *Event) Unmarshal(payload []byte) error {
	return json.Unmarshal(payload, e)
}

func (e *Event) GetGUID() []byte {
	return e.Guid
}

func (e *Event) GetTopic() string {
	return e.Topic
}
