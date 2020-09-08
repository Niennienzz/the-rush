package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"

	"the-rush/graph/model"
	"the-rush/repository"
)

func main() {
	var (
		ctx    = context.Background()
		repo   = repository.NewLocal()
		inputs []model.PlayerInput
	)

	data, err := ioutil.ReadFile("resources/rushing.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &inputs)
	if err != nil {
		panic(err)
	}

	for _, input := range inputs {
		player, err := input.ToPlayer()
		if err != nil {
			log.Println(err)
			continue
		}

		_, err = repo.Player().Create(ctx, player)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("ingested player name %s, with generated id %s\n", player.Name, player.ID)
	}
}
