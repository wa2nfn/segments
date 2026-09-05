package main

import "strings"

// Converts "38:148,98,48" to "rgb(148,98,48)"
func parseRGB(code string) string {
	parts := strings.Split(code, ":")
	if len(parts) != 2 {
		return ""
	}
	return "rgb(" + parts[1] + ")"
}
