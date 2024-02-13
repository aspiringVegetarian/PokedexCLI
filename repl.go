package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func startRepl(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)
	commandMap := loadCommands()

	for {
		fmt.Print("Pokedex >")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		if len(input) == 0 {
			continue
		}

		command, exists := commandMap[input[0]]
		if exists {
			if input[0] == "explore" {
				if len(input) == 1 {
					cfg.specificLocation = nil
				} else if len(input) == 2 {
					cfg.specificLocation = &input[1]
				} else if len(input) > 2 {
					fmt.Println("Only enter one location id or name after the explore command")
					continue
				}
			}
			if input[0] == "catch" || input[0] == "pokedex" {
				if len(input) == 1 {
					cfg.specificPokemon = nil
				} else if len(input) == 2 {
					cfg.specificPokemon = &input[1]
				} else if len(input) > 2 {
					if input[0] == "catch" {
						fmt.Println("Please enter only one Pokemon id or name after the catch command")
					} else {
						fmt.Println("Please enter only one Pokemon name after the pokedex command")
					}
					continue
				}
			}
			err := command.callback(cfg)
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}
