package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	mj "github.com/jbenet/go-multicodec/json"
	"github.com/kataras/iris"
	"github.com/kr/pretty"
	j "github.com/nimona/go-nimona/journal"
	"github.com/nimona/go-nimona/repository"
	"github.com/nimona/go-nimona/store"
)

func main() {
	f, err := os.OpenFile("/tmp/nimona-journal-kv.mjson", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	mc := mj.Codec(true)
	journal := j.NewJournal(mc, f, f)

	pairsRepositoryStore := store.NewInMemoryStore()
	pairsRepository := repository.NewRepository(pairsRepositoryStore, &KV{}, &Event{})
	journal.Notify(pairsRepository)

	journal.Rewind()

	api := &kvAPI{
		journal: journal,
		pairs:   pairsRepository,
	}

	iris.Get("/:key", api.Get)
	iris.Post("/:key", api.Set)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	iris.Listen(":" + port)
}

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

type KVSetEventPayload struct {
	Value *[]byte `json:"value"`
}

type KV struct {
	Key   []byte  `json:"k"`
	Value *[]byte `json:"v"`
}

func (kv *KV) GetGUID() []byte {
	return kv.Key
}

func (kv *KV) Marshal() ([]byte, error) {
	return json.Marshal(kv)
}

func (kv *KV) Unmarshal(v []byte) error {
	return json.Unmarshal(v, kv)
}

func (kv *KV) ApplyEvent(genericEvent repository.Event) {
	eventJSON, _ := genericEvent.Marshal() // TODO(geoah) Check error

	var payloadRawJSON json.RawMessage
	event := &Event{
		Payload: &payloadRawJSON,
	}
	json.Unmarshal(eventJSON, event) // TODO(geoah) Check error

	pretty.Println("> Applying", genericEvent, event.GetGUID(), event.GetTopic(), event.Payload)

	switch event.GetTopic() {
	case "set":
		payload := &KVSetEventPayload{}
		if err := json.Unmarshal(payloadRawJSON, payload); err != nil {
			log.Fatal(err)
			return
		}
		kv.Key = event.GetGUID()
		kv.Value = payload.Value
		if kv.Value != nil {
			fmt.Printf("* SET %s=%s\n", string(kv.Key), string(*kv.Value))
		} else {
			fmt.Printf("* DEL %s\n", string(kv.Key))
		}
	}
}

type kvAPI struct {
	journal *j.SerialJournal
	pairs   *repository.Repository
}

func (api *kvAPI) Get(c *iris.Context) {
	key := c.Param("key")
	pair, err := api.pairs.GetByGUID([]byte(key))
	if err != nil {
		c.Text(iris.StatusNotFound, "Not found")
		return
	}
	instanceConcrete, ok := pair.(*KV)
	if ok != true {
		c.Text(iris.StatusInternalServerError, "Could not cast") // TODO(geoah) Better error
		return
	}
	if instanceConcrete.Value != nil {
		c.Text(iris.StatusOK, string(*instanceConcrete.Value))
	} else {
		c.Text(iris.StatusNotFound, "Not found")
	}
}

func (api *kvAPI) Set(c *iris.Context) {
	key := c.Param("key")
	value := c.Request.Body()

	event := &Event{
		Guid:  []byte(key),
		Topic: "set",
		Payload: &KVSetEventPayload{
			Value: &value,
		},
	}

	eventJSON, _ := json.Marshal(event) // TODO(geoah) Handle error
	_, err := api.journal.Append(eventJSON)
	if err != nil {
		c.Text(iris.StatusInternalServerError, "Could not save value")
		log.Println("Could not append event. err=", err)
		return
	}
	instance, err := api.pairs.GetByGUID([]byte(key))
	if err != nil {
		c.Text(iris.StatusInternalServerError, "Could not get key")
		return
	}
	instanceConcrete, ok := instance.(*KV)
	if ok != true {
		c.Text(iris.StatusInternalServerError, "Could not cast") // TODO(geoah) Better error
		return
	}
	if instanceConcrete.Value != nil {
		c.Text(iris.StatusOK, string(*instanceConcrete.Value))
	} else {
		c.Text(iris.StatusNotFound, "Not found")
	}
}
