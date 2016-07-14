package main

import (
	"os"

	"github.com/kataras/iris"
	j "github.com/nimona/go-nimona/journal"
	"github.com/nimona/go-nimona/repository"
	"github.com/nimona/go-nimona/store"
)

func main() {
	journalStore := store.NewInMemoryStore()
	journal := j.NewJournal(journalStore)

	instanceRepositoryStore := store.NewInMemoryStore()
	instanceRepository := repository.NewInstanceRepository(journal, instanceRepositoryStore)
	journal.Notify(instanceRepository)

	api := NewAPI("1", journal, instanceRepository)

	iris.Get("/instances/:id", api.GetInstance)
	iris.Patch("/instances/:id", api.PatchInstance)
	iris.Post("/instances", api.PostInstance)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	iris.Listen(":" + port)
}
