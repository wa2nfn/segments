package main

import (
	"fmt"
	"math"
)

// toMixedNumber formats a decimal length in inches to a mixed number string,
func toMixedNumber(decimalInches float64) string {
	var fraction int = decimalPrecision
	var result string

	if decimalInches < 0 {
		decimalInches = math.Abs(decimalInches)
	}

	whole := math.Floor(decimalInches)
	fractional := decimalInches - whole

	numerator := int(math.Ceil(fractional * float64(fraction)))
	denominator := fraction

	if numerator == denominator {
		whole++
		result = fmt.Sprintf("%.0f", whole)
		return result
	}

	if numerator == 0 {
		if whole > 0 {
			result = fmt.Sprintf("%.0f", whole)
		} else {
			result = fmt.Sprintf("%.0f", float64(1/decimalPrecision))
		}
		return result
	}

	commonDivisor := gcd(numerator, denominator)

	simplifiedNum := numerator / commonDivisor
	simplifiedDen := denominator / commonDivisor

	// Build the final string.
	if whole > 0 {
		result = fmt.Sprintf("%.0f-%d/%d", whole, simplifiedNum, simplifiedDen)
	} else {
		result = fmt.Sprintf("%d/%d", simplifiedNum, simplifiedDen)
	}

	return result
}
