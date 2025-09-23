package repl

import (
	"fmt"
	"strconv"

	"github.com/cbrookscode/pokedexcli2/internal"
)

func DisplayLocations(locations internal.ListofLocations, config *Config) {
	config.Menu_options = Menu{Options: make(map[string]string)}
	fmt.Println("-------------------------------")
	for i, area := range locations.Results {
		fmt.Printf("%d) %v\n", i+1, area.Name)
		config.Menu_options.Options[strconv.Itoa(i+1)] = area.Name
	}
	fmt.Println("-------------------------------")
}

func DisplayPokedex(pokedex *internal.Pokedex, config *Config) {
	counter := 1
	config.Menu_options = Menu{Options: make(map[string]string)}
	fmt.Println("Your Pokedex:")
	for key := range pokedex.Entries {
		fmt.Printf("   %d) %v\n", counter, key)
		config.Menu_options.Options[strconv.Itoa(counter)] = key
		counter++
	}
	fmt.Println()
}

func DisplayPokemonInArea(area internal.LocationArea) {
	fmt.Println("---- Found Pokemon ----")
	for _, pokemon := range area.PokemonEncounters {
		fmt.Printf(" - %v\n", pokemon.Pokemon.Name)
	}
	fmt.Println("-----------------------")
}
