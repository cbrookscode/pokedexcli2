package main

import (
	"strings"
	"unicode"
)

func cleanInput(text string) []string {
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
