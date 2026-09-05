package main

// Escapes HTML-sensitive characters
func htmlEscape(s string) string {
	switch s {
	case "<":
		return "&lt;"
	case ">":
		return "&gt;"
	case "&":
		return "&amp;"
	default:
		return s
	}
}
