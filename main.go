package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/cbrookscode/pokedexcli2/internal"
	"github.com/cbrookscode/pokedexcli2/player"
	"github.com/cbrookscode/pokedexcli2/repl"
)

func main() {
	// get available commands for the program
	availableCommands := repl.RegisterCommands()

	cache := internal.NewCache(60 * time.Second)

	fd := os.Stdin.Fd()
	orig, err := repl.EnableRawMode(fd)
	if err != nil {
		fmt.Printf("error enabling raw mode for terminal input\n")
		os.Exit(1)
	}

	defer repl.DisableRawMode(fd, orig)

	term := &repl.Terminal{Line_prefix: "Pokedex >"}

	config := repl.Config{Next: "", Current: "", Previous: nil, Menu_options: repl.Menu{Options: make(map[int]string)}, Player: player.Player{Level: 5}}

	pokedex := internal.Pokedex{Entries: make(map[string]internal.Pokemon)}

	starting_pokemon_options := []string{"bulbasaur", "charmander", "squirtle"}

	start_of_program := true

	fmt.Println("Welcome! Please select a pokemon from the following options")
	for i, pokemon_name := range starting_pokemon_options {
		fmt.Printf("%d) %s\n", i+1, pokemon_name)
		config.Menu_options.Options[i+1] = pokemon_name
	}
	fmt.Println("Type the number on the left and press enter to make your selection")

	for {
		fmt.Printf("%s ", term.Line_prefix)
		repl.HandleUserInput(term, fd, orig)

		cleaned := repl.CleanInput(term.User_input)
		if len(cleaned) == 0 {
			continue
		}

		if start_of_program {
			num, err := strconv.Atoi(term.User_input) // check if provided input used a number based on display options
			if err == nil {
				_, found := config.Menu_options.Options[num] // confirm number provided aligns with current available options
				if !found {
					fmt.Println("the number you provided is outside the range of available displayed options")
					for i, pokemon_name := range starting_pokemon_options {
						fmt.Printf("%d) %s\n", i+1, pokemon_name)
					}
					continue
				}
				pokemon, err := internal.GetPokemon(config.Menu_options.Options[num]) // make get request for pokemon info
				if err != nil {
					fmt.Println("Error making get request for specified pokemon")
					return
				}
				pokedex.AddPokemonToPokedex(pokemon)
				err = config.Player.AddPokemonToPlayerParty(pokemon, &pokedex)
				if err != nil {
					fmt.Printf("%v\n", err)
				}
				fmt.Printf("%v has been added to your party\n", pokemon.Name)
				start_of_program = false
				continue
			} else {
				fmt.Println("Please type a number cooresponding to one of your available options and then press enter.")
				for i, pokemon_name := range starting_pokemon_options {
					fmt.Printf("%d) %s\n", i+1, pokemon_name)
				}
				continue
			}
		}

		// command for cli will always be first word in user input
		usersCommand := cleaned[0]

		command, ok := availableCommands[usersCommand]
		if !ok {
			fmt.Printf("Command provided(%v) does not exist\n", usersCommand)
			continue
		}
		err = command.Callback(&config, cache, &pokedex, term.User_input)
		if err != nil {
			if err.Error() == "exit" {
				break
			}
			fmt.Println(err)
			return
		}
	}
}
