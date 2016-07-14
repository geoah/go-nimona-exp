package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/kataras/iris"
	"github.com/nimona/go-nimona/journal"
	"github.com/nimona/go-nimona/repository"
	"github.com/pborman/uuid"
)

type API struct {
	ownerID   string
	journal   *journal.Journal
	instances *repository.InstanceRepository
}

func NewAPI(ownerID string, journal *journal.Journal, instanceRepository *repository.InstanceRepository) *API {
	return &API{
		ownerID:   ownerID,
		journal:   journal,
		instances: instanceRepository,
	}
}

func (api *API) GetInstance(c *iris.Context) {
	id := c.Param("id")
	instance, err := api.instances.GetResourceByID(id)
	if err != nil {
		c.JSON(iris.StatusNotFound, iris.Map{"error": "Not found"})
		return
	}
	c.JSON(iris.StatusOK, instance)
}

type InstancePostRequest struct {
	ID      string      `json:"id"`
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func (api *API) PostInstance(c *iris.Context) {
	instanceRequest := &InstancePostRequest{}
	err := c.ReadJSON(&instanceRequest)
	if err != nil {
		c.JSON(iris.StatusBadRequest, iris.Map{"error": "Invalid JSON"})
		return
	}

	id := uuid.New()
	instanceCreated := &repository.Event{
		Topic: "Created",
		Payload: &repository.InstanceCreatedEvent{
			ID:      id,
			Type:    instanceRequest.Type,
			OwnerID: api.ownerID,
			Created: time.Now(),
			Updated: time.Now(),
			Payload: instanceRequest.Payload,
		},
	}

	instanceCreatedJSON, _ := json.Marshal(instanceCreated) // TODO(geoah) Handle error
	err = api.journal.AppendPayload(instanceCreatedJSON)
	if err != nil {
		c.JSON(iris.StatusInternalServerError, iris.Map{"error": err})
		log.Println("Could not append payload. err=", err)
		return
	}
	instance, err := api.instances.GetResourceByID(id)
	if err != nil {
		c.JSON(iris.StatusInternalServerError, iris.Map{"error": err})
		return
	}
	c.JSON(iris.StatusOK, instance)
}

type InstancePatchRequest struct {
	Payload interface{} `json:"payload"`
}

func (api *API) PatchInstance(c *iris.Context) {
	instanceRequest := &InstancePatchRequest{}
	err := c.ReadJSON(&instanceRequest)
	if err != nil {
		c.JSON(iris.StatusBadRequest, iris.Map{"error": "Invalid JSON"})
		return
	}

	id := c.Param("id")
	instanceUpdated := &repository.Event{
		Topic: "Updated",
		Payload: &repository.InstanceUpdatedEvent{
			ID:      id,
			Updated: time.Now(),
			Payload: instanceRequest.Payload,
		},
	}

	instanceUpdatedJSON, _ := json.Marshal(instanceUpdated) // TODO(geoah) Handle error
	err = api.journal.AppendPayload(instanceUpdatedJSON)
	if err != nil {
		c.JSON(iris.StatusInternalServerError, iris.Map{"error": err})
		return
	}
	instance, err := api.instances.GetResourceByID(id)
	if err != nil {
		c.JSON(iris.StatusInternalServerError, iris.Map{"error": err})
		return
	}
	c.JSON(iris.StatusOK, instance)
}
