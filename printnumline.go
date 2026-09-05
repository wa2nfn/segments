package main

import (
	"fmt"
	"strings"
)

func printNumLine(DLimage string) (string, bool) {
	var rc bool
	runes := []rune(DLimage)
	chunkSize := 120
	var result strings.Builder

	if len(runes) > 120 {
		rc = true
	} else {
		chunkSize = len(runes)
	}

	// Loop through the string in chunks of 120 runes
	for start := 0; start < len(runes); start += chunkSize {
		end := start + chunkSize
		if end > len(runes) {
			end = len(runes)
		}

		// Add the double newline separator for subsequent chunks
		if start > 0 {
			result.WriteString("\n\n")
		}

		var line0, line1, line2 strings.Builder

		needsHundreds := end >= 100

		// Build the index ruler lines
		for i := start + 1; i <= end; i++ {
			if needsHundreds {
				if i < 100 {
					line0.WriteString(" ")
				} else {
					line0.WriteString(fmt.Sprintf("%d", (i/100)%10))
				}
			}

			if i < 10 {
				line1.WriteString(" ")
			} else {
				line1.WriteString(fmt.Sprintf("%d", (i/10)%10))
			}

			line2.WriteString(fmt.Sprintf("%d", i%10))
		}

		// Append the index lines (outputs 3 index rows if max index >= 100, otherwise 2)
		if needsHundreds {
			result.WriteString(line0.String() + "\n")
		}
		result.WriteString(line1.String() + "\n")
		result.WriteString(line2.String() + "\n")

		// Append the blank line
		result.WriteString("\n")

		// Append the actual text chunk
		textChunk := string(runes[start:end])
		result.WriteString(textChunk + "\n")

		// Append the colorized return string
		result.WriteString(colorize(textChunk))
	}

	pic := result.String()
	return pic, rc
}
