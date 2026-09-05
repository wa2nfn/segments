package main

import (
	"fmt"
	"strings"
)

// colorPattern iterates over the input string, applying True Color FG/BG escape sequences
// to each character and creating a block of color.
func colorPattern(input string, colorMap map[rune]TrueColor) string {
	var sb strings.Builder

	for _, char := range input {
		colors, ok := colorMap[char]
		if !ok {

			sb.WriteRune(char)
			sb.WriteString(ANSI_RESET)
			continue
		}

		fgEscape := fmt.Sprintf("\x1b[38;2;%d;%d;%dm", colors.R_FG, colors.G_FG, colors.B_FG)

		bgEscape := fmt.Sprintf("\x1b[48;2;%d;%d;%dm", colors.R_BG, colors.G_BG, colors.B_BG)

		sb.WriteString(bgEscape)

		sb.WriteString(fgEscape)

		sb.WriteRune(char)

		sb.WriteString(ANSI_RESET)
	}

	return sb.String()
}
