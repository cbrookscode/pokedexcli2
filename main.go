package main

import (
	"fmt"
	"os"
	"time"

	"github.com/cbrookscode/pokedexcli2/internal"
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

	term := repl.Terminal{
		History: []string{},
	}

	config := repl.Config{Next: "", Current: "", Previous: nil, Orig_Term_Settings: orig}

	pokedex := internal.Pokedex{Entries: make(map[string]internal.Pokemon)}

	for {
		repl.HandleUserInput(&term)

		cleaned := repl.CleanInput(term.User_input)
		if len(cleaned) == 0 {
			continue
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
			fmt.Println(err)
			return
		}
	}
}
