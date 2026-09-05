package main

import (
	"fmt"
	"regexp"
	"strings"
)

func colorize(image string) string {

	if colorType == "trueColor" {
		return colorPattern(image, colorMap)
	}

	reXs := regexp.MustCompile(`D+`)
	reSpaces := regexp.MustCompile(`L+`)
	reOther := regexp.MustCompile(`O+`)
	coloredImage := image

	coloredImage = reXs.ReplaceAllStringFunc(coloredImage, func(s string) string {
		spaces := strings.Repeat("D", len(s))
		return fmt.Sprintf("%s%s%s", ANSI_BROWN, spaces, ANSI_RESET)
	})

	coloredImage = reSpaces.ReplaceAllStringFunc(coloredImage, func(s string) string {
		spaces := strings.Repeat("L", len(s))
		return fmt.Sprintf("%s%s%s", ANSI_BEIGE, spaces, ANSI_RESET)
	})

	coloredImage = reOther.ReplaceAllStringFunc(coloredImage, func(s string) string {
		spaces := strings.Repeat("O", len(s))
		return fmt.Sprintf("%s%s%s", ANSI_TAN, spaces, ANSI_RESET)
	})

	coloredImage = reOther.ReplaceAllStringFunc(coloredImage, func(s string) string {
		spaces := strings.Repeat("L", len(s))
		return fmt.Sprintf("%s%s%s", ANSI_TAN, spaces, ANSI_RESET)
	})

	coloredImage = strings.ReplaceAll(coloredImage, "D", " ")
	coloredImage = strings.ReplaceAll(coloredImage, "O", " ")
	coloredImage = strings.ReplaceAll(coloredImage, "L", " ")
	return coloredImage
}
