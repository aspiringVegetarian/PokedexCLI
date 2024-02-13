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
				"         Each call of the command displays the names of 20 Pokemon in the Pokemon world.\n" +
				"         Each subsequent call to the command will display the next 20 pokemon.",
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
		"pokedex": {
			name: "pokedex",
			description: "Show information for any Pokemon in your Pokedex (must have encountered via catch command).\n" +
				"         Provide a Pokemon name following the command.\n" +
				"         If you do not provide a name following the command, all of the Pokemon in your Pokedex will be listed.",
			callback: commandPokedex,
		},
		"team": {
			name:        "team",
			description: "Lists the Pokemon you have caught.",
			callback:    commandTeam,
		},
	}
}

func commandHelp(cfg *config) error {

	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range loadCommands() {
		fmt.Printf("%s: %s\n\n", cmd.name, cmd.description)
	}

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
	}
	if cfg.specificPokemon == nil {
		randomID := rand.Intn(*cfg.pokemonCount-1) + 1
		if randomID > 1025 {
			randomID = 10000 + (randomID - 1025)
		}
		randomIDString := fmt.Sprint(randomID)
		cfg.specificPokemon = &randomIDString

	}

	resp, err := cfg.pokeapiClient.ExplorePokemon(cfg.specificPokemon)
	if err != nil {
		return err
	}

	fmt.Printf("\nYou see a wild %s!\n", resp.Name)

	givenName, caught := cfg.pokedexCaught[resp.Name]
	if caught {
		fmt.Printf("\nYou have already caught a %s and named it %s, so you let this one go...\n\n", resp.Name, givenName)
		return nil
	}

	_, seen := cfg.pokedexSeen[resp.Name]
	if seen {
		fmt.Printf("\nIt won't get away this time!\n")
	}

	fmt.Printf("\nYou throw a PokeBall at %s...\n", resp.Name)

	catchDifficulty := (resp.BaseExperience / 2) * (resp.Height / 4) * (resp.Weight / 4)
	catchProb := 0.90
	if catchDifficulty > 100000 {
		catchProb = 0.10
	} else if catchDifficulty > 75000 {
		catchProb = 0.20
	} else if catchDifficulty > 50000 {
		catchProb = 0.30
	} else if catchDifficulty > 25000 {
		catchProb = 0.50
	} else if catchDifficulty > 10000 {
		catchProb = 0.70
	} else if catchDifficulty > 5000 {
		catchProb = 0.80
	}

	if rand.Float64() < catchProb {
		// CAUGHT
		fmt.Printf("\n%s was successfully caught!\n\n", resp.Name)
		fmt.Println("Name your newly caught Pokemon: ")
		var newName string
		fmt.Scanln(&newName)
		if newName == "" {
			newName = resp.Name
		}
		cfg.pokedexCaught[resp.Name] = newName
		if !seen {
			fmt.Printf("\n%s has been added to your team and %s has been added to your Pokedex!\n\n", newName, resp.Name)
		} else {
			fmt.Printf("\n%s was added to your team!\n\n", newName)
		}
		cfg.pokedexSeen[resp.Name] = resp

	} else {
		// GOT AWAY
		fmt.Printf("\n%s got away!\n", resp.Name)
		if !seen {
			fmt.Printf("\nHowever, %s has been added to your Pokedex!\n\n", resp.Name)
			cfg.pokedexSeen[resp.Name] = resp
		} else {
			fmt.Println()
		}
	}

	return nil
}

func commandPokedex(cfg *config) error {

	if cfg.specificPokemon == nil {
		//k := rand.Intn(len(cfg.pokedexSeen))
		fmt.Printf("\nYou have the following Pokemon in your Pokedex: ")
		for name, _ := range cfg.pokedexSeen {
			fmt.Printf("\n * %s", name)
		}
		fmt.Printf("\n\nUse the pokedex command with any of the Pokemon names listed to see more info.\n")
		return nil
	} else {
		info := cfg.pokedexSeen[*cfg.specificPokemon]
		givenName, caught := cfg.pokedexCaught[info.Name]
		if caught {
			if givenName != info.Name {
				fmt.Printf("\nYou have caught a %s and named it %s!\n", info.Name, givenName)
			} else {
				fmt.Printf("\nYou have caught a %s!\n", info.Name)
			}
		}
		fmt.Printf("\nName: %s", info.Name)
		fmt.Printf("\nHeight: %v", info.Height)
		fmt.Printf("\nWeight: %v", info.Weight)
		fmt.Printf("\nStats:")
		for _, content := range info.Stats {
			fmt.Printf("\n  --%s: %v", content.Stat.Name, content.BaseStat)
		}
		fmt.Printf("\nType(s):")
		for _, content := range info.Types {
			fmt.Printf("\n  --%s", content.Type.Name)
		}
		fmt.Printf("\n\n")

		return nil

	}
}

func commandTeam(cfg *config) error {
	if len(cfg.pokedexCaught) == 0 {
		fmt.Printf("\nYou haven't caught any Pokemon yet! Get out there!\n\n")
		return nil
	}
	fmt.Printf("\nYou have caught the following Pokemon: ")
	for name, givenName := range cfg.pokedexCaught {
		if givenName != name {
			fmt.Printf("\n * %s the %s", givenName, name)
		} else {
			fmt.Printf("\n * %s", name)
		}
	}
	fmt.Printf("\n\nThey are your team, treat them well!\n\n")

	return nil
}
