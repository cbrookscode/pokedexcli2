package repl

import (
	"fmt"
	"os"

	"github.com/cbrookscode/pokedexcli2/internal"
)

type cliCommand struct {
	name        string
	description string
	Callback    func() error
}

type Config struct {
	Next     string
	Previous any
}

var pagination = Config{
	Next:     "",
	Previous: nil,
}

func RegisterCommands() map[string]cliCommand {
	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			Callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Provide a help menu to explain options for user",
			Callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Get the next 20 location areas",
			Callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Get the previous 20 location areas",
			Callback:    commandMapb,
		},
	}
	return commands
}
func commandExit() error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("=================================")
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for key, command := range RegisterCommands() {
		fmt.Printf("%v: %v\n", key, command.description)
	}
	fmt.Println("=================================")
	return nil
}

func commandMap() error {
	url := "https://pokeapi.co/api/v2/location-area/"
	if pagination.Next != "" {
		url = pagination.Next
	}

	locations, err := internal.GetLocations(url)
	if err != nil {
		return fmt.Errorf("error using map command: %v", err)
	}
	pagination.Next = locations.Next
	pagination.Previous = locations.Previous

	for _, area := range locations.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func commandMapb() error {
	url, ok := pagination.Previous.(string)
	if !ok {
		fmt.Println("there is no previous page to go back to")
		return nil
	}

	locations, err := internal.GetLocations(url)
	if err != nil {
		return fmt.Errorf("error getting locations: %v", err)
	}

	pagination.Next = locations.Next
	pagination.Previous = locations.Previous

	for _, area := range locations.Results {
		fmt.Println(area.Name)
	}
	return nil
}
