package repl

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/cbrookscode/pokedexcli2/internal"
)

func FilterInputForMenuOptionSelection(cleaned_user_input []string, config *Config) (string, error) {
	num, err := strconv.Atoi(cleaned_user_input[1]) // check if provided input used a number based on display options
	if err == nil {
		_, found := config.Menu_options.Options[num] // confirm number provided aligns with current available options
		if !found {
			return "", fmt.Errorf("the number you provided is outside the range of available displayed options")
		}
		return config.Menu_options.Options[num], nil
	}
	return cleaned_user_input[1], nil
}

func GrabPokemon(cleaned_user_input []string, config *Config, pokedex *internal.Pokedex) (internal.Pokemon, error) {
	var pokemon internal.Pokemon

	entered_pokemon_name, err := FilterInputForMenuOptionSelection(cleaned_user_input, config)
	if err != nil {
		return internal.Pokemon{}, err
	}
	pokemon, err = pokedex.GetPokemonFromPokedex(entered_pokemon_name) // grab pokemon info from pokedex if available
	if err != nil {
		pokemon, err = internal.GetPokemon(entered_pokemon_name) // make get request for pokemon
		if err != nil {
			return internal.Pokemon{}, fmt.Errorf("issue grabbing requested pokemon")
		}
	}
	return pokemon, nil
}

func ThrowPokeball(pokemon internal.Pokemon, pokedex *internal.Pokedex) {
	fmt.Printf("Throwing a Pokeball at %v...\n", pokemon.Name)

	// Calculate chance to catch and determine if successful. Add to pokedex on success
	myRNG := rand.Intn(101)
	difficulty := internal.CalcChancetoCatchDifficulty(pokemon.BaseExperience)
	if float64(myRNG) >= difficulty {
		pokedex.AddPokemonToPokedex(pokemon)
		fmt.Printf("You've caught %v!\n\n", pokemon.Name)
		return
	}
	fmt.Printf("%v managed to break free!\n\n", pokemon.Name)
}
