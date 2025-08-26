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
	input_bytes []byte
	index       int
	history     []string
	User_input  string
	cursor      int
}

func redrawTerminal(term *Terminal, calc_backstep bool) {
	fmt.Print("\r\033[K")
	fmt.Printf("Pokedex > %s", term.input_bytes)
	if calc_backstep {
		num_backsteps := len(term.input_bytes) - term.cursor
		if num_backsteps > 0 {
			fmt.Printf("\033[%dD", num_backsteps)
		}
	} else {
		term.cursor = len(term.input_bytes)
	}
}

func HandleUserInput(term *Terminal) {
	exit := false
	term.cursor = 0
	fmt.Print("Pokedex > ")
	for {
		buf := make([]byte, 1)
		_, err := os.Stdin.Read(buf)
		if err != nil {
			fmt.Printf("error reading user input\n")
			os.Exit(2)
		}
		switch buf[0] {
		case '\r', '\n': // enter
			fmt.Print("\n")
			term.history = append(term.history, string(term.input_bytes))
			term.index = len(term.history)
			term.User_input = string(term.input_bytes)
			term.input_bytes = term.input_bytes[:0]
			exit = true
		case '\b', '\x7F': // backspace
			if len(term.input_bytes) == 0 {
				continue
			}
			if term.cursor > 0 {
				term.input_bytes = append(term.input_bytes[:term.cursor-1], term.input_bytes[term.cursor:]...)
				term.cursor--
				redrawTerminal(term, true)
			}
		case '\x1b': // escape key / start of arrow keys
			_, err = os.Stdin.Read(buf)
			if err != nil {
				fmt.Printf("error reading user input\n")
				os.Exit(2)
			}
			switch buf[0] {
			case '[':
				os.Stdin.Read(buf)
				switch buf[0] {
				case 'A': // up arrow
					if len(term.history) < 1 {
						continue
					}
					if term.index-1 >= 0 {
						term.index--
						term.input_bytes = []byte(term.history[term.index])
						redrawTerminal(term, false)
					}
				case 'B': // down arrow
					if len(term.history) < 1 {
						continue
					}
					if term.index+1 >= len(term.history) {
						continue
					}
					term.index++
					term.input_bytes = []byte(term.history[term.index])
					redrawTerminal(term, false)
				case 'C': // right arrow
					if term.cursor < len(term.input_bytes) {
						fmt.Print("\033[C")
						term.cursor++
					}
				case 'D': // left arrow
					if term.cursor > 0 {
						fmt.Print("\033[D")
						term.cursor--
					}
				}
			}
		default:
			if term.cursor >= 0 && term.cursor <= len(term.input_bytes) {
				term.input_bytes = append(term.input_bytes[:term.cursor], append([]byte{buf[0]}, term.input_bytes[term.cursor:]...)...)
				term.cursor++
				redrawTerminal(term, true)
			}
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
