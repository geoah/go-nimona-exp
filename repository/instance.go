package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/nimona/go-nimona/journal"
	"github.com/nimona/go-nimona/store"
)

// NewClusteringKey creates a new ClusteringKey using the user's ID.
func NewClusteringKey(resourceID string) store.ClusteringKey {
	return &ClusteringKey{
		keys: []store.Key{
			resourceID,
		},
	}
}

type Event struct {
	Topic   string      `json:"topic"`
	Payload interface{} `json:"payload"`
}

type EventUnknown struct {
	Topic   string          `json:"topic"`
	Payload json.RawMessage `json:"payload"`
}

type InstanceCreatedEvent struct {
	ID      string      `json:"id"`
	OwnerID string      `json:"owner_id"`
	Type    string      `json:"type"`
	Created time.Time   `json:"created"`
	Updated time.Time   `json:"updated"`
	Payload interface{} `json:"payload"`
}

type InstanceUpdatedEvent struct {
	ID      string      `json:"id"`
	Updated time.Time   `json:"updated"`
	Payload interface{} `json:"payload"`
}

type InstanceRemovedEvent struct {
	ID      string    `json:"id"`
	Updated time.Time `json:"updated"`
}

type Instance struct {
	ID      string      `json:"id"`
	OwnerID string      `json:"owner_id"`
	Type    string      `json:"type"`
	Created time.Time   `json:"created"`
	Updated time.Time   `json:"updated"`
	Removed bool        `json:"removed"`
	Payload interface{} `json:"payload"`
}

func NewInstance(id, ownerID, itype string, payload interface{}) *Instance {
	return &Instance{
		ID:      id,
		OwnerID: ownerID,
		Type:    itype,
		Payload: payload,
		Created: time.Now(),
		Updated: time.Now(),
	}
}

func (i *Instance) ToJSON() ([]byte, error) {
	return json.Marshal(i)
}

type InstanceRepository struct {
	journal *journal.SerialJournal
	store   store.Store
}

func NewInstanceRepository(j *journal.SerialJournal, s store.Store) *InstanceRepository {
	return &InstanceRepository{
		journal: j,
		store:   s,
	}
}

func (r *InstanceRepository) GetResourceByID(id string) (*Instance, error) {
	key := NewClusteringKey(id)
	instanceJSON, err := r.store.GetOne(key)
	if err != nil {
		return nil, err
	}
	instance := &Instance{}
	errUnmarshalling := json.Unmarshal(instanceJSON, instance)
	if errUnmarshalling != nil {
		return nil, errUnmarshalling
	}
	return instance, nil
}

func (r *InstanceRepository) AppendEntry(entry journal.Entry) (Resource, error) {
	return &Instance{}, nil
}

func (r *InstanceRepository) AppendedEntry(entry journal.Entry) {
	fmt.Println("> Processing", string(entry.GetPayload()))
	event := &EventUnknown{}
	err := json.Unmarshal(entry.GetPayload(), event)
	if err != nil {
		// TODO(geoah) Log error.
		return
	}

	switch event.Topic {
	case "Created":
		fmt.Println(">> As created")
		eventPayload := &InstanceCreatedEvent{}
		errPayload := json.Unmarshal(event.Payload, eventPayload)
		if errPayload == nil {
			r.handleEvent(eventPayload)
		}
	case "Updated":
		fmt.Println(">> As Updated")
		eventPayload := &InstanceUpdatedEvent{}
		errPayload := json.Unmarshal(event.Payload, eventPayload)
		if errPayload == nil {
			r.handleEvent(eventPayload)
		}
	case "Removed":
		fmt.Println(">> As removed")
		eventPayload := &InstanceRemovedEvent{}
		errPayload := json.Unmarshal(event.Payload, eventPayload)
		if errPayload == nil {
			r.handleEvent(eventPayload)
		}
	}
}

func (r *InstanceRepository) handleEvent(event interface{}) error {
	// TODO(geoah) Refactor.
	fmt.Println("> Handling event")
	switch t := event.(type) {
	case *InstanceCreatedEvent:
		fmt.Println(">> As created")
		instance := &Instance{
			ID:      t.ID,
			OwnerID: t.OwnerID,
			Type:    t.Type,
			Created: t.Created,
			Updated: t.Updated,
			Payload: t.Payload,
		}
		key := NewClusteringKey(t.ID)
		instanceJSON, _ := json.Marshal(instance)
		err := r.store.Put(key, instanceJSON) // TODO(geoah) Handle error
		return err
	case *InstanceUpdatedEvent:
		fmt.Println(">> As Updated", t)
		instance, err := r.GetResourceByID(t.ID)
		if err != nil {
			return err
		}
		instance.Updated = t.Updated
		instance.Payload = t.Payload
		key := NewClusteringKey(t.ID)
		instanceJSON, _ := json.Marshal(instance)
		err = r.store.Put(key, instanceJSON) // TODO(geoah) Handle error
		return err
	case *InstanceRemovedEvent:
		fmt.Println(">> As removed")
		instance, err := r.GetResourceByID(t.ID)
		if err != nil {
			return err
		}
		instance.Updated = t.Updated
		instance.Removed = true
		key := NewClusteringKey(t.ID)
		instanceJSON, _ := json.Marshal(instance)
		err = r.store.Put(key, instanceJSON) // TODO(geoah) Handle error
		return err
	default:
		return errors.New("erm... invalid event")
	}
	return nil
}
