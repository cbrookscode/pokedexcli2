package main

import (
	"fmt"
	"os"

	"github.com/cbrookscode/pokedexcli2/repl"
)

func main() {
	// get available commands for the program
	// availableCommands := repl.RegisterCommands()

	// cache := internal.NewCache(60 * time.Second)

	fd := os.Stdin.Fd()
	orig, err := repl.EnableRawMode(fd)
	if err != nil {
		fmt.Printf("error enabling raw mode for terminal input\n")
		os.Exit(1)
	}
	defer repl.DisableRawMode(fd, orig)

	// config := repl.Config{Next: "", Current: "", Previous: nil, Orig_Term_Settings: orig}

	// pokedex := internal.Pokedex{Entries: make(map[string]internal.Pokemon)}

	// listen for user input
	// scanner := bufio.NewScanner(os.Stdin)

	history := []string{}
	count := 1
	input_bytes := []byte{}
	fmt.Print("Pokedex > ")
	for {

		buf := make([]byte, 1)

		_, err := os.Stdin.Read(buf)
		if err != nil {
			fmt.Printf("error reading user input\n")
			os.Exit(2)
		}
		switch buf[0] {
		case '\r', '\n':
			fmt.Print("\n")
			fmt.Print("Pokedex > ")
			history = append(history, string(input_bytes))
			count = 1
			input_bytes = input_bytes[:0]
		case '\b', '\x7F':
			if len(input_bytes) == 0 {
				continue
			}
			input_bytes = input_bytes[:len(input_bytes)-1]
			fmt.Print("\b \b")
		case '\x1b':
			_, err = os.Stdin.Read(buf)
			if err != nil {
				fmt.Printf("error reading user input\n")
				os.Exit(2)
			}
			switch buf[0] {
			case '[':
				os.Stdin.Read(buf)
				switch buf[0] {
				case 'A':
					if len(history) < 1 || len(history)-count < 0 {
						continue
					}
					fmt.Print("\r\033[2K")
					fmt.Print("Pokedex > ")
					previous_input := history[len(history)-count]
					fmt.Printf("%s", previous_input)
					input_bytes = []byte(previous_input)
					count++
				case 'B':
					count--
					if len(history) < 1 || len(history)-count > len(history) {
						count++
						continue
					}
					fmt.Print("\r\033[2K")
					fmt.Print("Pokedex > ")
					previous_input := history[len(history)-count]
					fmt.Printf("%s", previous_input)
					input_bytes = []byte(previous_input)
				case 'C':
				case 'D':
				}
			}
		default:
			input_bytes = append(input_bytes, buf[0])
			fmt.Printf("%c", buf[0])
		}
	}
}

// advance scanner to next token
// ready := scanner.Scan()
// if !ready {
// 	fmt.Printf("Error scanning for user input: %v", scanner.Err())
// 	continue
// }

// grab user input from scanner in string format
// input := scanner.Text()
// cleaned := repl.CleanInput(input)
// if len(cleaned) == 0 {
// 	fmt.Print("No command provided\n")
// 	continue
// }
// command for cli will always be first word in user input
// usersCommand := cleaned[0]

// command, ok := availableCommands[usersCommand]
// if !ok {
// 	fmt.Printf("Command provided(%v) does not exist\n\n", usersCommand)
// 	continue
// }
// err := command.Callback(&config, cache, &pokedex, input)
// if err != nil {
// 	fmt.Println(err)
// 	return
// }
