package main

import (
	"fmt"
	"log"
	"pokedex/api"
)

func main() {
	pokemon, err := api.GetPokemon(25)
	if err != nil {
		log.Fatalf("Impossible de récuperer le pokémon", err)
	}
	fmt.Println(pokemon)
}
