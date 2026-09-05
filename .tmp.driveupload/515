package main

import (
	"fmt"
	"strings"
)

// display morse and patterns for each character
func showMorse(MorseTable map[rune]Morse, length int) {

	ClearScreen()
	screen.Printf("                                   Morse Code Characters/Patterns\n\n")

	screen.Printf("  Name   %-8s  %-23s  %-23s  %s  %s  %s\n", "Morse", "D&L Segments", "Color Segments", "# Dark", "# Light", "Total")
	screen.Printf("  ----   %8s  %-23s  %-23s  %-s  %-s  %-s\n",
		strings.Repeat("-", 8), strings.Repeat("-", 23), strings.Repeat("-", 23), strings.Repeat("-", 6), strings.Repeat("-", 7), strings.Repeat("-", 6))

	order := "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890.,?/=+-:;\"'@* "

	for _, keyRune := range order {
		pattern := MorseTable[keyRune]

		if _, ok := MorseTable[keyRune]; !ok {
			fmt.Printf("Missing entry for rune: %q (%d)\n", keyRune, keyRune)
		}

		pattern.Cimage = colorize(pattern.DLimage)
		segSum := pattern.Light + pattern.Dark

		padCount := 23 - len(pattern.DLimage)
		padding := strings.Repeat(" ", padCount)

		screen.Printf("  '%c'   %-9s  %-23s  %-s%s     %2d       %2d      %2d\n\n",
			keyRune, pattern.Morse, pattern.DLimage, pattern.Cimage, padding, pattern.Dark, pattern.Light, segSum)
	}
}
