package repl

import (
	"strings"
	"syscall"
	"unicode"
	"unsafe"
)

// passed fd can be obtained by os.stdin.fd(). returned termious struct pointer will be original terminal settings to revert back to
func enableRawMode(fd int) (*syscall.Termios, error) {
	// struct containing settings for terminal
	orig := &syscall.Termios{}

	// system call to linux kernal to obtain terminal settings. first variable is type of system call, next is getting file descriptor for stdin which is the location we want to be working in,
	// tcgets is getting the terminal settings, last arguement is struct to put info into.
	_, _, err := syscall.Syscall6(uintptr(syscall.SYS_IOCTL), uintptr(fd), uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(orig)), 0, 0, 0)
	if err != 0 {
		return nil, err
	}

	// dereference and store copy of original settings
	raw := *orig

	// turn off canonical & echo. | operator combines the two bits into bitmask and then &^= will clear any bit field in lflag that is turned on in the bitmask created by combining echo and icanon.
	// icanon controls whether or not terminal holds it as a buffer until user hits enter, echo controls printing to terminal.
	raw.Lflag &^= syscall.ICANON | syscall.ECHO

	_, _, err = syscall.Syscall6(uintptr(syscall.SYS_IOCTL), uintptr(fd), uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&raw)), 0, 0, 0)
	if err != 0 {
		return nil, err
	}

	return orig, nil
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
