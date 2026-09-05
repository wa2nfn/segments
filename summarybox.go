package main

import (
	"fmt"
	"strings"
)

// BuildSummaryBox creates a 4-line box around S1 and S2.
// Title is centered, S1 and S2 are padded, and all widths adjust automatically.
func BuildSummaryBox(title, S1, S2, S3 string) string {
	// Determine the widest content line (S1 is always >= S2 per your note)
	contentWidth := len(S2)
	if S2 == "" {
		contentWidth = len(S1)
	}

	// The inside width of the box is contentWidth + 2 spaces (one after "= " and one before " =")
	innerWidth := contentWidth + 4

	// Center the title inside the top line
	titleLine := centerText(title, innerWidth)

	// Build the top and bottom border of '='
	border := strings.Repeat("=", len(titleLine))

	// Build the S1 and S2 lines
	line1 := fmt.Sprintf("= %-*s =", contentWidth, S1)
	var line2 string
	if S2 != "" {
		line2 = fmt.Sprintf("= %-*s =", contentWidth, S2) // left-align S2, pad to S1 width
	}
	line3 := fmt.Sprintf("= %-*s =", contentWidth, S3) // left-align S2, pad to S1 width

	// Assemble final box
	var res string
	if S2 != "" {
		res = fmt.Sprintf("%s\n%s\n%s\n%s\n%s",
			border,
			titleLine,
			line1,
			line2,
			line3,
		) + "\n" + border
	} else {
		res = fmt.Sprintf("%s\n%s\n%s\n%s",
			border,
			titleLine,
			line1,
			line3,
		) + "\n" + border
	}

	return res
}

// centerText centers text inside a field of width w using '=' padding.
func centerText(text string, w int) string {
	if len(text) >= w {
		return text
	}
	totalPad := w - len(text)
	left := totalPad / 2
	right := totalPad - left
	return strings.Repeat("=", left) + text + strings.Repeat("=", right)
}
