package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/cbrookscode/pokedexcli2/repl"
)

func main() {
	// get available commands for the program
	availableCommands := repl.RegisterCommands()

	// listen for user input
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		// advance scanner to next token
		ready := scanner.Scan()
		if !ready {
			fmt.Printf("Error scanning for user input: %v", scanner.Err())
			continue
		}

		// grab user input from scanner in string format
		input := scanner.Text()
		cleaned := repl.CleanInput(input)
		if len(cleaned) == 0 {
			fmt.Print("No command provided\n")
			continue
		}
		// command for cli will always be first word in user input
		usersCommand := cleaned[0]

		command, ok := availableCommands[usersCommand]
		if !ok {
			fmt.Printf("Command provided(%v) does not exist\n", usersCommand)
			continue
		}
		err := command.Callback()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
