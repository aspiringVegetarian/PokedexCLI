package main

import (
	"time"

	"github.com/aspiringVegetarian/PokedexCLI/internal/pokeapi"
)

type config struct {
	pokeapiClient       pokeapi.Client
	nextLocationAreaURL *string
	prevLocationAreaURL *string
	locationCount       *int
	specificLocation    *string
	nextPokemonURL      *string
	prevPokemonURL      *string
	pokemonCount        *int
	specificPokemon     *string
}

func main() {
	cfg := config{
		pokeapiClient: pokeapi.NewClient(time.Minute),
	}
	startRepl(&cfg)
}
