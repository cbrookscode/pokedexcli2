package repl

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unicode"
	"unsafe"
)

type Terminal struct {
	Input_bytes []byte
	Index       int
	History     []string
	User_input  string
}

func HandleUserInput(term *Terminal) {
	exit := false

	fmt.Print("Pokedex > ")
	for {
		buf := make([]byte, 1)

		_, err := os.Stdin.Read(buf)
		if err != nil {
			fmt.Printf("error reading user input\n")
			os.Exit(2)
		}
		switch buf[0] {
		// enter
		case '\r', '\n':
			fmt.Print("\n")
			term.History = append(term.History, string(term.Input_bytes))
			term.Index = len(term.History)
			term.User_input = string(term.Input_bytes)
			term.Input_bytes = term.Input_bytes[:0]
			exit = true
		// backspace or del
		case '\b', '\x7F':
			if len(term.Input_bytes) == 0 {
				continue
			}
			term.Input_bytes = term.Input_bytes[:len(term.Input_bytes)-1]
			fmt.Print("\b \b")
		// escape key / start of arrow keys
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
					if len(term.History) < 1 {
						continue
					}
					if term.Index-1 >= 0 {
						term.Index--
						fmt.Print("\r\033[2K")
						fmt.Print("Pokedex > ")
						previous_input := term.History[term.Index]
						fmt.Printf("%s", previous_input)
						term.Input_bytes = []byte(previous_input)
					}
				case 'B':
					if len(term.History) < 1 {
						continue
					}
					if term.Index+1 >= len(term.History) {
						continue
					}
					term.Index++
					fmt.Print("\r\033[2K")
					fmt.Print("Pokedex > ")
					previous_input := term.History[term.Index]
					fmt.Printf("%s", previous_input)
					term.Input_bytes = []byte(previous_input)
				case 'C':
				case 'D':
				}
			}
		default:
			term.Input_bytes = append(term.Input_bytes, buf[0])
			fmt.Printf("%c", buf[0])
		}
		if exit {
			break
		}
	}
}

// passed fd can be obtained by os.stdin.fd(). returned termious struct pointer will be original terminal settings to revert back to
func EnableRawMode(fd uintptr) (*syscall.Termios, error) {
	// struct containing settings for terminal
	orig := &syscall.Termios{}

	// system call to linux kernal to obtain terminal settings. first variable is type of system call, next is getting file descriptor for stdin which is the location we want to be working in,
	// tcgets is getting the terminal settings, last arguement is struct to put info into.
	_, _, err := syscall.Syscall6(uintptr(syscall.SYS_IOCTL), fd, uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(orig)), 0, 0, 0)
	if err != 0 {
		return nil, err
	}

	// dereference and store copy of original settings
	raw := *orig

	// turn off canonical & echo. | operator combines the two bits into bitmask and then &^= will clear any bit field in lflag that is turned on in the bitmask created by combining echo and icanon.
	// icanon controls whether or not terminal holds it as a buffer until user hits enter, echo controls printing to terminal.
	raw.Lflag &^= syscall.ICANON | syscall.ECHO

	_, _, err = syscall.Syscall6(uintptr(syscall.SYS_IOCTL), fd, uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&raw)), 0, 0, 0)
	if err != 0 {
		return nil, err
	}

	return orig, nil
}

func DisableRawMode(fd uintptr, orig *syscall.Termios) error {
	_, _, err := syscall.Syscall6(uintptr(syscall.SYS_IOCTL), fd, uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(orig)), 0, 0, 0)
	if err != 0 {
		return err
	}

	return nil
}

func CleanInput(text string) []string {
	cleaned := []string{}
	for _, word := range strings.Fields(strings.ToLower(text)) {
		newword := ""
		for _, letter := range word {
			if unicode.IsLetter(letter) {
				newword = newword + string(letter)
			}
		}
		cleaned = append(cleaned, newword)
	}
	return cleaned
}
