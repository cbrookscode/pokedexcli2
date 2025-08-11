package repl

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cbrookscode/pokedexcli2/internal"
)

type cliCommand struct {
	name        string
	description string
	Callback    func(config *Config, cache *internal.Cache) error
}

type Config struct {
	Next     string
	Previous any
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
func commandExit(config *Config, cache *internal.Cache) error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config, cache *internal.Cache) error {
	fmt.Println("=================================")
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for key, command := range RegisterCommands() {
		fmt.Printf("%v: %v\n", key, command.description)
	}
	fmt.Println("=================================")
	return nil
}

// need to make check against default start point in case the program reaches the end of next pages available. Next would also be "" in this case.
func commandMap(config *Config, cache *internal.Cache) error {
	if config.Next == "" {
		url := "https://pokeapi.co/api/v2/location-area/"
		val, ok := cache.Get(url)
		if ok {
			locations := internal.ListofLocations{}
			err := json.Unmarshal(val, &locations)
			if err != nil {
				return fmt.Errorf("error unmarshalling data from cache into list of locations struct: %v", err)
			}
			fmt.Println("*** got information from cache ***")
			for _, area := range locations.Results {
				fmt.Println(area.Name)
			}

			config.Previous = nil
			config.Next = locations.Next

			return nil
		}

		locations, bytes, err := internal.GetLocations(url)
		if err != nil {
			return fmt.Errorf("error using map command: %v", err)
		}

		for _, area := range locations.Results {
			fmt.Println(area.Name)
		}

		config.Next = locations.Next
		config.Previous = nil
		cache.Add(url, bytes)

		return nil
	}

	url := config.Next
	fmt.Printf("url to get is: %v\n", url)
	val, ok := cache.Get(url)
	// found cached entry for next page request
	if ok {
		locations := internal.ListofLocations{}
		err := json.Unmarshal(val, &locations)
		if err != nil {
			return fmt.Errorf("error unmarshalling data from cache into list of locations struct: %v", err)
		}
		fmt.Println("*** got information from cache ***")
		for _, area := range locations.Results {
			fmt.Println(area.Name)
		}

		config.Next = locations.Next
		config.Previous = locations.Previous

		return nil
	}

	// make new get request, and cache data
	locations, bytes, err := internal.GetLocations(url)
	fmt.Printf("url to get is: %v\n", url)
	if err != nil {
		return fmt.Errorf("error using map command: %v", err)
	}

	for _, area := range locations.Results {
		fmt.Println(area.Name)
	}

	config.Next = locations.Next
	config.Previous = locations.Previous
	cache.Add(url, bytes)
	fmt.Println("cache added for recent url")

	return nil
}

func commandMapb(config *Config, cache *internal.Cache) error {
	url, ok := config.Previous.(string)
	if !ok {
		fmt.Println("there is no previous page to go back to")
		return nil
	}

	// check if cached
	val, ok := cache.Get(url)
	if ok {
		locations := internal.ListofLocations{}
		err := json.Unmarshal(val, &locations)
		if err != nil {
			return fmt.Errorf("error unmarshalling data from cache into list of locations struct: %v", err)
		}
		fmt.Println("*** got information from cache ***")
		for _, area := range locations.Results {
			fmt.Println(area.Name)
		}

		config.Next = locations.Next
		config.Previous = locations.Previous

		return nil
	}

	// not cached, make new request and cache data
	locations, bytes, err := internal.GetLocations(url)
	if err != nil {
		return fmt.Errorf("error getting locations: %v", err)
	}

	for _, area := range locations.Results {
		fmt.Println(area.Name)
	}

	config.Next = locations.Next
	config.Previous = locations.Previous
	cache.Add(url, bytes)

	return nil
}
