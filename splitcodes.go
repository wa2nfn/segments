package main

import "strings"

// Splits ANSI codes into usable chunks
func splitCodes(codes string) []string {
	parts := strings.Split(codes, ";")
	var result []string
	for i := 0; i < len(parts); {
		if parts[i] == "38" || parts[i] == "48" {
			if i+4 < len(parts) && parts[i+1] == "2" {
				result = append(result, parts[i]+":"+parts[i+2]+","+parts[i+3]+","+parts[i+4])
				i += 5
				continue
			}
		}
		result = append(result, parts[i])
		i++
	}
	return result
}
