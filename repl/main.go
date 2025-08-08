package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	availableCommands := registerCommands()

	for {
		fmt.Print("Pokedex > ")

		ready := scanner.Scan()
		if !ready {
			fmt.Printf("Error scanning for user input: %v", scanner.Err())
			continue
		}

		input := scanner.Text()
		cleaned := cleanInput(input)
		if len(cleaned) == 0 {
			fmt.Print("No command provided\n")
			continue
		}
		usersCommand := cleaned[0]

		command, ok := availableCommands[usersCommand]
		if !ok {
			fmt.Printf("Command provided(%v) does not exist\n", usersCommand)
			continue
		}
		command.callback()
	}
}
