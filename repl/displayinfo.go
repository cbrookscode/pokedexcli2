package repl

import (
	"fmt"

	"github.com/cbrookscode/pokedexcli2/internal"
)

func DisplayLocations(locations internal.ListofLocations, config *Config) {
	config.Menu_options = Menu{Options: make(map[int]string)}
	fmt.Println("-------------------------------")
	for i, area := range locations.Results {
		fmt.Printf("%d) %v\n", i+1, area.Name)
		config.Menu_options.Options[i+1] = area.Name
	}
	fmt.Println("-------------------------------")
}

func DisplayPokedex(pokedex *internal.Pokedex, config *Config) {
	counter := 1
	config.Menu_options = Menu{Options: make(map[int]string)}
	fmt.Println("Your Pokedex:")
	for key := range pokedex.Entries {
		fmt.Printf("   %d) %v\n", counter, key)
		config.Menu_options.Options[counter] = key
		counter++
	}
	fmt.Println()
}

func DisplayPokemonInArea(area internal.LocationArea, config *Config) {
	fmt.Println("---- Found Pokemon ----")
	for i, pokemon := range area.PokemonEncounters {
		fmt.Printf(" %d) %v\n", i+1, pokemon.Pokemon.Name)
		config.Menu_options.Options[i+1] = pokemon.Pokemon.Name
	}
	fmt.Println("-----------------------")
}

func DisplayPokemonInfoFromPokedex(pokemon internal.Pokemon) {
	// Print Pokemon info
	fmt.Printf("Name: %v\n", pokemon.Name)
	fmt.Println("Stats:")
	for _, statstruct := range pokemon.Stats {
		fmt.Printf("   -%v: %v\n", statstruct.Stat.Name, statstruct.BaseStat)
	}
	fmt.Println("Types:")
	for _, types := range pokemon.Types {
		fmt.Printf("   - %v\n", types.Type.Name)
	}
	fmt.Println("Moves:")
	for _, move := range pokemon.Moves {
		fmt.Printf("   - %v\n", move.Move.Name)
	}
}

func DisplayPokemonInfoFromParty(pokemon internal.Pokemon) {
	fmt.Printf("Name: %v\n", pokemon.Name)
	fmt.Println("Current Stats:")
	fmt.Printf(
		"   Attack: %v\n   Defense: %v\n   Max Hp: %v\n   Current Hp: %v\n   Special Attack: %v\n   Special Defense: %v\n   Speed: %v\n",
		pokemon.Current_stats.Attack,
		pokemon.Current_stats.Defense,
		pokemon.Current_stats.Hp,
		pokemon.Current_health,
		pokemon.Current_stats.Special_attack,
		pokemon.Current_stats.Special_defense,
		pokemon.Current_stats.Speed,
	)
	fmt.Println("Moves:")
	for _, move := range pokemon.Moves {
		fmt.Printf("   - %v\n", move.Move.Name)
	}
}
