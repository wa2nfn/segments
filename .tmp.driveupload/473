package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// use wedges file instead of built-ins
// note input can set color, precision or segments of wedges
// the processing will also skip any vlaues commented out with # or //
// Processing can be terminated with: stop, end, quit, about
func readWedgesFile() {
	const filename = "wedges.txt"
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	fmt.Printf("\nFYI: using your %s file.\n\n", filename)

	// Clear existing wedges
	for k := range stdWedges {
		delete(stdWedges, k)
	}

	var num = 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		num++
		line := strings.TrimSpace(scanner.Text())
		line = strings.ToUpper(line)

		if line == "" {
			continue
		}

		//  skip comment lines and stop processing on 4 cmds
		if strings.HasPrefix(line, "#") {
			continue
		} else if strings.HasPrefix(line, "//") {
			continue
		} else if strings.HasPrefix(line, "end") {
			fmt.Println("got here")
			return
		} else if strings.HasPrefix(line, "stop") {
			return
		} else if strings.HasPrefix(line, "quit") {
			return
		} else if strings.HasPrefix(line, "abort") {
			return
		}

		// --- NEW PRECISION HANDLING ---
		if strings.HasPrefix(line, "PRECISION=") {
			val := strings.TrimPrefix(line, "PRECISION=")

			switch val {
			case "16":
				decimalPrecision = 16
			case "32":
				decimalPrecision = 32
			case "64":
				decimalPrecision = 64
			default:
				// silently default to 32
				decimalPrecision = 32
			}
			continue
		}

		// --- NEW KERF HANDLING ---
		if strings.HasPrefix(line, "KERF=") {

			line = strings.TrimPrefix(line, "KERF=")
			if strings.Contains(line, "/") {
				kerf, _ = parseFractionalFloat(line)
			} else {
				kerf, _ = strconv.ParseFloat(line, 64)
			}
			continue
		}

		// --- NEW COLOR HANDLING ---
		if strings.HasPrefix(line, "COLOR=") {
			val := strings.TrimPrefix(line, "COLOR=")

			switch val {
			case "RED":
				updateColor('D', TrueColor{
					R_FG: 220, G_FG: 20, B_FG: 60,
					R_BG: 220, G_BG: 20, B_BG: 60,
				})
				updateColor('O', TrueColor{
					R_FG: 255, G_FG: 140, B_FG: 0,
					R_BG: 255, G_BG: 140, B_BG: 0,
				})
				updateColor('L', TrueColor{
					R_FG: 255, G_FG: 215, B_FG: 0,
					R_BG: 255, G_BG: 215, B_BG: 0,
				})
			case "BLUE":
				updateColor('D', TrueColor{
					R_FG: 30, G_FG: 30, B_FG: 180,
					R_BG: 30, G_BG: 30, B_BG: 180,
				})
				updateColor('O', TrueColor{
					R_FG: 70, G_FG: 130, B_FG: 230,
					R_BG: 70, G_BG: 130, B_BG: 230,
				})
				updateColor('L', TrueColor{
					R_FG: 0, G_FG: 180, B_FG: 255,
					R_BG: 0, G_BG: 180, B_BG: 255,
				})
			case "GREEN":
				updateColor('D', TrueColor{
					R_FG: 0, G_FG: 100, B_FG: 0,
					R_BG: 0, G_BG: 100, B_BG: 0,
				})
				updateColor('O', TrueColor{
					R_FG: 34, G_FG: 139, B_FG: 34,
					R_BG: 34, G_BG: 139, B_BG: 34,
				})
				updateColor('L', TrueColor{
					R_FG: 0, G_FG: 255, B_FG: 127,
					R_BG: 0, G_BG: 255, B_BG: 127,
				})
			default:
				// stays BROWN
			}
			continue
		}

		// --- EXISTING WEDGE HANDLING ---
		val, err := strconv.Atoi(line)
		if err != nil {
			fmt.Printf("warning: not an integer <%q> ignoring line number <%d>\n", line, num)
			continue
		}
		if val < 4 || val > 180 {
			continue
		}
		if _, exists := stdWedges[val]; exists {
			fmt.Printf("warning: duplicate wedge: %d\n", val)
		} else {
			stdWedges[val] = struct{}{}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error reading file: %v\n", err)
	}
}

// Function to update a key
func updateColor(key rune, newColor TrueColor) {
	colorMap[key] = newColor
}
