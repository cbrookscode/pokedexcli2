package repl

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"github.com/cbrookscode/pokedexcli2/internal"
)

type cliCommand struct {
	name        string
	description string
	Callback    func(config *Config, cache *internal.Cache, pokedex *internal.Pokedex, input string) error
}

type Menu struct {
	Options map[int]string
}

type Config struct {
	Next               string
	Current            string
	Previous           any
	Orig_Term_Settings *syscall.Termios
	File_desc          uintptr
	Menu_options       Menu
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
		"inspect": {
			name:        "inspect",
			description: "Get information from your Pokedex about the provided pokemon",
			Callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "See what pokemon you have captured",
			Callback:    commandPokedex,
		},
	}
	return commands
}

func commandExit(config *Config, cache *internal.Cache, pokedex *internal.Pokedex, input string) error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	return errors.New("exit")
}

func commandHelp(config *Config, cache *internal.Cache, pokedex *internal.Pokedex, input string) error {
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
func commandMap(config *Config, cache *internal.Cache, pokedex *internal.Pokedex, input string) error {
	direction := "forward"
	if config.Next == "" {
		url := "https://pokeapi.co/api/v2/location-area/"
		// check if cached
		bytes, ok := cache.Get(url)
		if ok {
			locations, err := internal.GetLocationsFromCache(bytes)
			if err != nil {
				return err
			}
			DisplayLocations(locations, config)
			err = updatePreviousAndNext(&locations, config, direction, url)
			if err != nil {
				return err
			}
			return nil
		}

		// not cached make new request
		locations, bytes, err := internal.GetLocations(url)
		if err != nil {
			return err
		}
		DisplayLocations(locations, config)
		err = updatePreviousAndNext(&locations, config, direction, url)
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
		locations, err := internal.GetLocationsFromCache(bytes)
		if err != nil {
			return err
		}
		DisplayLocations(locations, config)
		err = updatePreviousAndNext(&locations, config, direction, url)
		if err != nil {
			return err
		}
		return nil
	}

	// not cached make new request
	locations, bytes, err := internal.GetLocations(url)
	if err != nil {
		return err
	}
	DisplayLocations(locations, config)
	err = updatePreviousAndNext(&locations, config, direction, url)
	if err != nil {
		return err
	}
	cache.Add(url, bytes)

	return nil
}

func commandMapb(config *Config, cache *internal.Cache, pokedex *internal.Pokedex, input string) error {
	direction := "backward"
	url, ok := config.Previous.(string)
	if !ok || url == "" {
		fmt.Printf("there is no previous page to go back to\n\n")
		return nil
	}

	// check if cached
	bytes, ok := cache.Get(url)
	if ok {
		locations, err := internal.GetLocationsFromCache(bytes)
		if err != nil {
			return err
		}
		DisplayLocations(locations, config)
		err = updatePreviousAndNext(&locations, config, direction, url)
		if err != nil {
			return err
		}
		return nil
	}

	// not cached, make new request
	locations, bytes, err := internal.GetLocations(url)
	if err != nil {
		return err
	}
	DisplayLocations(locations, config)
	err = updatePreviousAndNext(&locations, config, direction, url)
	if err != nil {
		return err
	}
	cache.Add(url, bytes)

	return nil
}

// provide list of pokemon in provided location
func commandExplore(config *Config, cache *internal.Cache, pokedex *internal.Pokedex, input string) error {
	cleanedInput := strings.Fields(strings.ToLower(input))
	if len(cleanedInput) != 2 {
		fmt.Println("please provide a location name after explore or a number matching to the currently displayed areas via map or mapb command. Do not provide further info. Ex: explore canalave-city-area, explore 10")
		return nil
	}

	num, err := strconv.Atoi(cleanedInput[1])
	if err == nil {
		_, found := config.Menu_options.Options[num]
		if !found {
			fmt.Println("The number you provided is outside the range of available displayed options.")
			return nil
		}
		area, err := internal.GetLocationsPokemon(config.Menu_options.Options[num])
		if err != nil {
			return fmt.Errorf("issue grabbing this locations available pokemon: %w", err)
		}
		DisplayPokemonInArea(area, config)
		return nil
	}

	area, err := internal.GetLocationsPokemon(cleanedInput[1])
	if err != nil {
		fmt.Println("issue grabbing this location. Location name is incorrectly spelled or doesn't exist")
		return nil
	}
	DisplayPokemonInArea(area, config)
	return nil
}

func commandCatch(config *Config, cache *internal.Cache, pokedex *internal.Pokedex, input string) error {
	// get pokemon name
	cleanedInput := strings.Fields(strings.ToLower(input))
	if len(cleanedInput) != 2 {
		fmt.Println("please provide a pokemon name after catch or a number that refers to displayed options. Do not provide any further info. Ex: catch pikachu, catch 2")
		return nil
	}

	pokemon := internal.Pokemon{}

	num, err := strconv.Atoi(cleanedInput[1])
	if err == nil {
		_, found := config.Menu_options.Options[num]
		if !found {
			fmt.Println("The number you provided is outside the range of available displayed options.")
			return nil
		}
		// grab pokemon info if in pokedex already
		pokemon, err = pokedex.GetPokemonFromPokedex(config.Menu_options.Options[num])
		if err != nil {
			// make get request for pokemon xp
			pokemon, err = internal.GetPokemon(config.Menu_options.Options[num])
			if err != nil {
				fmt.Printf("Pokemon to grab: %s", config.Menu_options.Options[num])
				fmt.Println("issue grabbing this pokemon")
				return nil
			}
		}
	} else {
		// grab pokemon info if in pokedex already
		pokemon, err = pokedex.GetPokemonFromPokedex(cleanedInput[1])
		if err != nil {
			// make get request for pokemon xp
			pokemon, err = internal.GetPokemon(cleanedInput[1])
			if err != nil {
				fmt.Println("issue grabbing this pokemon. Pokemon name is incorrectly spelled or doesn't exist")
				return nil
			}
		}
	}

	fmt.Printf("Throwing a Pokeball at %v...\n", pokemon.Name)

	// Calculate chance to catch and determine if successful. Add to pokedex
	myRNG := rand.Intn(101)
	difficulty := internal.CalcChancetoCatchDifficulty(pokemon.BaseExperience)
	if float64(myRNG) >= difficulty {
		pokedex.AddPokemonToPokedex(pokemon)
		fmt.Printf("You've caught %v!\n\n", pokemon.Name)
		return nil
	}
	fmt.Printf("%v managed to break free!\n", pokemon.Name)

	fmt.Println()
	return nil
}

func commandInspect(config *Config, cache *internal.Cache, pokedex *internal.Pokedex, input string) error {
	// get pokemon name
	cleanedInput := strings.Fields(strings.ToLower(input))
	if len(cleanedInput) != 2 {
		fmt.Println("please provide a pokemon name after inspect or a number related to the most recent display. Do not provide further info. Ex: inspect pikachu, inspect 2")
		return nil
	}

	num, err := strconv.Atoi(cleanedInput[1])
	if err == nil {
		_, found := config.Menu_options.Options[num]
		if !found {
			fmt.Println("The number you provided is outside the range of available displayed options.")
			return nil
		}
		pokemon, err := pokedex.GetPokemonFromPokedex(config.Menu_options.Options[num])
		if err != nil {
			fmt.Printf("The pokemon name provided doesn't exist in your pokedex or was spell incorrectly\n\n")
			return nil
		}
		DisplayPokemonInfo(pokemon)
		return nil
	}

	pokemon, err := pokedex.GetPokemonFromPokedex(cleanedInput[1])
	if err != nil {
		fmt.Printf("The pokemon name provided doesn't exist in your pokedex or was spell incorrectly\n\n")
		return nil
	}
	DisplayPokemonInfo(pokemon)
	return nil
}

func commandPokedex(config *Config, cache *internal.Cache, pokedex *internal.Pokedex, input string) error {
	DisplayPokedex(pokedex, config)
	return nil
}
