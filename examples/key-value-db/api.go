package main

import (
	"encoding/json"
	"log"

	"github.com/kataras/iris"
	j "github.com/nimona/go-nimona/journal"
	"github.com/nimona/go-nimona/repository"
)

type Api struct {
	journal *j.SequentialJournal
	pairs   *repository.Repository
}

func (api *Api) Get(c *iris.Context) {
	key := c.Param("key")
	pair, err := api.pairs.GetByGUID([]byte(key))
	if err != nil {
		c.Text(iris.StatusNotFound, "Not found")
		return
	}
	instanceConcrete, ok := pair.(*Pair)
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

func (api *Api) Set(c *iris.Context) {
	key := c.Param("key")
	value := c.Request.Body()

	event := &Event{
		Guid:  []byte(key),
		Topic: "set",
		Payload: &PairSetEventPayload{
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
	instanceConcrete, ok := instance.(*Pair)
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
