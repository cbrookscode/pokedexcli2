package repl

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"

	"github.com/cbrookscode/pokedexcli2/internal"
)

type cliCommand struct {
	name        string
	description string
	Callback    func(config *Config, cache *internal.Cache, input string) error
}

type Config struct {
	Next     string
	Current  string
	Previous any
}

// To facilitate the ability to loop when moving back and forward in map and mapb commands. If this logic doesn't exist then when you get to first entry then you cant go back to last.
func updatePreviousAndNext(locations *internal.ListofLocations, config *Config, direction string, url string) error {
	switch direction {
	case "forward":
		config.Next = locations.Next
		config.Previous = config.Current
		locations.Previous = config.Current
		config.Current = url
	case "backward":
		config.Next = config.Current
		locations.Next = config.Current
		config.Previous = locations.Previous
		config.Current = url
	default:
		return errors.New("provided unexpected direction")
	}
	return nil
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
		"explore": {
			name:        "mapb",
			description: "See what pokemon are available for a given location area",
			Callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Try to catch a pokemon",
			Callback:    commandCatch,
		},
	}
	return commands
}

func commandExit(config *Config, cache *internal.Cache, input string) error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config, cache *internal.Cache, input string) error {
	// grab command list to sort key names into ordered list for better help menu ui experience
	commands := RegisterCommands()
	orderedKeys := []string{}
	for key := range commands {
		orderedKeys = append(orderedKeys, key)
	}
	sort.Strings(orderedKeys)

	fmt.Println("=================================")
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")

	for _, key := range orderedKeys {
		fmt.Printf("%v: %v\n", key, commands[key].description)
	}
	fmt.Println("=================================")
	return nil
}

// need to make check against default start point in case the program reaches the end of next pages available. Next would also be "" in this case.
func commandMap(config *Config, cache *internal.Cache, input string) error {
	direction := "forward"
	if config.Next == "" {
		url := "https://pokeapi.co/api/v2/location-area/"
		// check if cached
		bytes, ok := cache.Get(url)
		if ok {
			locations, err := internal.PrintLocationsFromCache(bytes)
			if err != nil {
				return err
			}
			err = updatePreviousAndNext(&locations, config, direction, url)
			if err != nil {
				return err
			}
			return nil
		}

		// not cached make new request
		locations, err := internal.PrintLocations(url)
		if err != nil {
			return err
		}
		err = updatePreviousAndNext(&locations, config, direction, url)
		if err != nil {
			return err
		}

		// add to cache
		bytes, err = json.Marshal(locations)
		if err != nil {
			return err
		}
		cache.Add(url, bytes)

		return nil
	}

	url := config.Next
	// check if cached
	bytes, ok := cache.Get(url)
	if ok {
		locations, err := internal.PrintLocationsFromCache(bytes)
		if err != nil {
			return err
		}
		err = updatePreviousAndNext(&locations, config, direction, url)
		if err != nil {
			return err
		}
		return nil
	}

	// not cached make new request
	locations, err := internal.PrintLocations(url)
	if err != nil {
		return err
	}
	err = updatePreviousAndNext(&locations, config, direction, url)
	if err != nil {
		return err
	}

	// add to cache
	bytes, err = json.Marshal(locations)
	if err != nil {
		return err
	}
	cache.Add(url, bytes)

	return nil
}

func commandMapb(config *Config, cache *internal.Cache, input string) error {
	direction := "backward"
	url, ok := config.Previous.(string)
	if !ok || url == "" {
		fmt.Println("there is no previous page to go back to")
		return nil
	}

	// check if cached
	bytes, ok := cache.Get(url)
	if ok {
		locations, err := internal.PrintLocationsFromCache(bytes)
		if err != nil {
			return err
		}
		err = updatePreviousAndNext(&locations, config, direction, url)
		if err != nil {
			return err
		}
		return nil
	}

	// not cached, make new request
	locations, err := internal.PrintLocations(url)
	if err != nil {
		return err
	}
	err = updatePreviousAndNext(&locations, config, direction, url)
	if err != nil {
		return err
	}

	// add to cache
	bytes, err = json.Marshal(locations)
	if err != nil {
		return err
	}
	cache.Add(url, bytes)

	return nil
}

// provide list of pokemon in provided location
func commandExplore(config *Config, cache *internal.Cache, input string) error {
	cleanedInput := strings.Fields(strings.ToLower(input))
	if len(cleanedInput) < 2 {
		fmt.Println("please provide a location name after explore. Ex: explore canalave-city-area")
		return nil
	}

	area, err := internal.GetLocationsPokemon(cleanedInput[1])
	if err != nil {
		fmt.Println("issue grabbing this location. Location name is incorrectly spelled or doesn't exist")
		return nil
	}

	fmt.Println("---- Found Pokemon ----")
	for _, pokemon := range area.PokemonEncounters {
		fmt.Printf(" - %v\n", pokemon.Pokemon.Name)
	}
	fmt.Println("-----------------------")

	return nil
}

func commandCatch(config *Config, cache *internal.Cache, input string) error {
	// get pokemon name
	cleanedInput := strings.Fields(strings.ToLower(input))
	if len(cleanedInput) < 2 {
		fmt.Println("please provide a pokemon name after catch. Ex: catch pikachu")
		return nil
	}

	// make get request for pokemon xp
	pokemon, err := internal.GetPokemon(cleanedInput[1])
	if err != nil {
		fmt.Println("issue grabbing this pokemon. Pokemon name is incorrectly spelled or doesn't exist")
		return nil
	}

	// get random number between 1 - 100
	myRNG := rand.Intn(101)

	// if calcchancetocatch number greater or equal to random number, then its caught, else it fails.
	difficulty := internal.CalcChancetoCatchDifficulty(pokemon.BaseExperience)
	if float64(myRNG) >= difficulty {
		fmt.Printf("You've caught %v!\n", pokemon.Name)
		return nil
	}
	fmt.Printf("%v managed to break free!\n", pokemon.Name)
	return nil
}
