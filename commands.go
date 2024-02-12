package main

import (
	"fmt"
	"math/rand"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func loadCommands() map[string]cliCommand {

	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name: "map",
			description: "Explore locations that can be visited within the Pokemon games.\n" +
				"     Each call of the command displays the names of 20 location areas in the Pokemon world.\n" +
				"     Each subsequent call to the command will display the next 20 locations.",
			callback: commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Explore the last 20 locations shown by the map command.",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore a specific location. Pass in a valid location name or id following the command, or it will select a random location.",
			callback:    commandExplore,
		},
		"pokemon": {
			name: "pokemon",
			description: "Explore Pokemon that can be caught within the Pokemon games.\n" +
				"     Each call of the command displays the names of 20 Pokemon in the Pokemon world.\n" +
				"     Each subsequent call to the command will display the next 20 pokemon.",
			callback: commandPokemon,
		},
		"pokemonb": {
			name:        "pokemonb",
			description: "Explore the last 20 Pokemon shown by the pokemon command.",
			callback:    commandPokemonb,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a specific Pokemon. Pass in a valid Pokemon name or id following the command, or it will try to catch a random Pokemon.",
			callback:    commandCatch,
		},
	}
}

func commandHelp(cfg *config) error {

	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range loadCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println()

	return nil
}

func commandExit(cfg *config) error {

	fmt.Println("Thank you for using PokedexCLI! See you soon")
	os.Exit(0)
	return nil
}

func commandMap(cfg *config) error {

	resp, err := cfg.pokeapiClient.ListLocationAreas(cfg.nextLocationAreaURL)
	if err != nil {
		return err
	}

	fmt.Println("Location areas:")
	for _, area := range resp.Results {
		fmt.Printf(" * %s\n", area.Name)
	}

	cfg.nextLocationAreaURL = resp.Next
	cfg.prevLocationAreaURL = resp.Previous
	if cfg.locationCount == nil {
		cfg.locationCount = resp.Count
	}

	return nil
}

func commandMapb(cfg *config) error {

	if cfg.prevLocationAreaURL == nil {
		return fmt.Errorf("You are on the first page. Call map again before using mapb (map back).")
	}

	resp, err := cfg.pokeapiClient.ListLocationAreas(cfg.prevLocationAreaURL)
	if err != nil {
		return err
	}

	fmt.Println("Location areas:")
	for _, area := range resp.Results {
		fmt.Printf(" * %s\n", area.Name)
	}

	cfg.nextLocationAreaURL = resp.Next
	cfg.prevLocationAreaURL = resp.Previous
	if cfg.locationCount == nil {
		cfg.locationCount = resp.Count
	}

	return nil
}

func commandExplore(cfg *config) error {
	if cfg.locationCount == nil && cfg.specificLocation == nil {
		resp, err := cfg.pokeapiClient.ListLocationAreas(cfg.nextLocationAreaURL)
		if err != nil {
			return err
		}
		cfg.locationCount = resp.Count
	}
	if cfg.specificLocation == nil {
		randomIDString := fmt.Sprint(rand.Intn(*cfg.locationCount-1) + 1)
		cfg.specificLocation = &randomIDString

	}

	resp, err := cfg.pokeapiClient.ExploreLocationArea(cfg.specificLocation)
	if err != nil {
		return err
	}

	fmt.Printf("\nExploring %s...\n\n", resp.Name)
	if len(resp.PokemonEncounters) == 0 {
		fmt.Println("There are no Pokemon here.")
	} else {
		fmt.Println("Found the following Pokemon: ")
		for _, encounter := range resp.PokemonEncounters {
			fmt.Printf(" * %s\n", encounter.Pokemon.Name)
		}
	}
	return nil
}
func commandPokemon(cfg *config) error {

	resp, err := cfg.pokeapiClient.ListPokemon(cfg.nextPokemonURL)
	if err != nil {
		return err
	}

	fmt.Println("Pokemon:")
	for _, area := range resp.Results {
		fmt.Printf(" * %s\n", area.Name)
	}

	cfg.nextPokemonURL = resp.Next
	cfg.prevPokemonURL = resp.Previous
	if cfg.pokemonCount == nil {
		cfg.pokemonCount = resp.Count
	}

	return nil
}

func commandPokemonb(cfg *config) error {

	if cfg.prevPokemonURL == nil {
		return fmt.Errorf("You are on the first page. Call pokemon again before using pokemonb (pokemon back).")
	}

	resp, err := cfg.pokeapiClient.ListPokemon(cfg.prevPokemonURL)
	if err != nil {
		return err
	}

	fmt.Println("Pokemon:")
	for _, area := range resp.Results {
		fmt.Printf(" * %s\n", area.Name)
	}

	cfg.nextPokemonURL = resp.Next
	cfg.prevPokemonURL = resp.Previous
	if cfg.pokemonCount == nil {
		cfg.pokemonCount = resp.Count
	}

	return nil
}

func commandCatch(cfg *config) error {
	if cfg.pokemonCount == nil && cfg.specificPokemon == nil {
		resp, err := cfg.pokeapiClient.ListPokemon(cfg.nextPokemonURL)
		if err != nil {
			return err
		}
		cfg.pokemonCount = resp.Count
		fmt.Println(*cfg.pokemonCount)
	}
	if cfg.specificPokemon == nil {
		randomID := rand.Intn(*cfg.pokemonCount-1) + 1
		if randomID > 1025 {
			randomID = 10000 + (randomID - 1025)
		}
		randomIDString := fmt.Sprint(randomID)
		fmt.Println(randomIDString)
		cfg.specificPokemon = &randomIDString

	}

	resp, err := cfg.pokeapiClient.ExplorePokemon(cfg.specificPokemon)
	if err != nil {
		return err
	}

	fmt.Printf("\nAttempting to catch %s...\n\n", resp.Name)

	return nil
}
