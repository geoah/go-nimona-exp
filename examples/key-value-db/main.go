package main

import (
	"os"

	mj "github.com/jbenet/go-multicodec/json"
	"github.com/kataras/iris"
	j "github.com/nimona/go-nimona/journal"
	"github.com/nimona/go-nimona/repository"
	"github.com/nimona/go-nimona/store"
)

func main() {
	// Create a file for our journal,
	f, err := os.OpenFile("/tmp/nimona-journal-kv.mjson", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// create a JSON multicodec to en/decode our journal enties in our file,
	mc := mj.Codec(true)
	// and create a journal.
	journal := j.NewJournal(mc, f, f)

	// Create an in-memory store for the repository,
	pairsRepositoryStore := store.NewInMemoryStore()
	// create a new repository; it requires a store, an aggregation target, and the events
	// that will be aggregated.
	// Since in our key-value we have only one one Event and one Aggregation it's pretty simple.
	// The repository will go through the events, and for each new `Event.GetGUID()` it gets will
	// create a new Aggregate and `aggregate.Apply(event)` on its events.
	pairsRepository := repository.NewRepository(pairsRepositoryStore, &Pair{}, &Event{})
	// ask to recieve notifications of new entries being added in the journal
	journal.Notify(pairsRepository)

	// we can now replay all entries in the journal to re-construct the state of our key-value store
	journal.Replay()

	// initialize our kv api
	api := &Api{
		journal: journal,
		pairs:   pairsRepository,
	}

	// setup routes
	iris.Get("/:key", api.Get)
	iris.Post("/:key", api.Set)

	// find out the port to listen to
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// and start the http server
	iris.Listen(":" + port)
}
