package main

import (
	"fmt"
	"os"
	"syscall"
	"unicode"

	"golang.org/x/term"
)

func readMessageCooked(
	MorseTable map[rune]Morse,
	wordSpaces int,
	interCharSpaces int,
	itsMorse bool,
) (visible string, cookedWord string, err error) {

	fd := int(syscall.Stdin)

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", "", fmt.Errorf("MakeRaw: %w", err)
	}
	defer term.Restore(fd, oldState)

	var (
		lineRunes []rune
		charSums  []int
		cooked    []rune
		lineTotal int
	)

	redraw := func() {
		text := string(lineRunes)
		var pending int
		if itsMorse {
			pending = wordSpaces
		} else {
			pending = interCharSpaces
		}
		fmt.Printf("\r\"%s\"  <%d> (+%d) Post message spaces.  ",
			text, lineTotal, pending)
	}

	buf := make([]byte, 1)

	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			return "", "", fmt.Errorf("read: %w", err)
		}
		if n == 0 {
			continue
		}

		b := buf[0]

		switch b {

		// ---------------------------------------------------------
		// RETURN — finalize
		// ---------------------------------------------------------
		case '\r', '\n':
			if len(lineRunes) == 0 {
				return "", "", nil
			}

			// Add final gap only if last rune is not a user space
			if len(lineRunes) > 0 && lineRunes[len(lineRunes)-1] != ' ' {
				var gap int
				if itsMorse {
					gap = wordSpaces
				} else {
					gap = interCharSpaces
				}

				lineTotal += gap
				charSums = append(charSums, gap)
				cooked = append(cooked, '#')

				if lineTotal > MAX_WORD_SEGMENTS {
					lineTotal -= gap
					charSums = charSums[:len(charSums)-1]
					cooked = cooked[:len(cooked)-1]
					fmt.Print("\r** TOO LONG — BACKSPACE TO FIX **")
					continue
				}
			}

			// Trim trailing user spaces (zero cost)
			for len(lineRunes) > 0 && lineRunes[len(lineRunes)-1] == ' ' {
				lineRunes = lineRunes[:len(lineRunes)-1]
				cooked = cooked[:len(cooked)-1]
			}

			// Trim trailing inter-character gap
			if len(cooked) > 0 && cooked[len(cooked)-1] == '#' {
				cooked = cooked[:len(cooked)-1]
			}

			return string(lineRunes), string(cooked), nil

		// ---------------------------------------------------------
		// CTRL+C
		// ---------------------------------------------------------
		case 3:
			fmt.Print("\n")
			return "", "", fmt.Errorf("")

			// ---------------------------------------------------------
			// BACKSPACE — fully corrected
			// ---------------------------------------------------------
		case 8, 127:
			if len(lineRunes) == 0 {
				// nothing to delete
				continue
			}

			// 1) If the last cooked entry is a timing marker (# or %),
			//    remove its timing from charSums and lineTotal.
			if len(cooked) > 0 {
				lastCooked := cooked[len(cooked)-1]
				if lastCooked == '#' || lastCooked == '%' {
					// gap timing is always the last charSums entry
					lineTotal -= charSums[len(charSums)-1]
					charSums = charSums[:len(charSums)-1]
					cooked = cooked[:len(cooked)-1]
				}
			}

			// 2) Now remove the letter itself: last cooked rune and last charSums entry.
			if len(cooked) > 0 && len(charSums) > 0 {
				lineTotal -= charSums[len(charSums)-1]
				charSums = charSums[:len(charSums)-1]
				cooked = cooked[:len(cooked)-1]
			}

			// 3) Remove the visible letter from lineRunes.
			lineRunes = lineRunes[:len(lineRunes)-1]

			redraw()
			continue

		// ---------------------------------------------------------
		// NORMAL CHARACTERS
		// ---------------------------------------------------------
		default:
			if b < 32 {
				continue
			}

			r := rune(b)

			if r == '#' || r == '%' {
				continue
			}

			if !(unicode.IsLetter(r) || unicode.IsDigit(r) || isAllowedPunc(r) || r == ' ') {
				continue
			}

			if r >= 'a' && r <= 'z' {
				r = r - 'a' + 'A'
			}

			// -----------------------------------------------------
			// USER SPACE → '%'
			// -----------------------------------------------------
			if r == ' ' {
				if len(lineRunes) == 0 {
					continue
				}
				if lineRunes[len(lineRunes)-1] == ' ' {
					continue
				}

				lineRunes = append(lineRunes, ' ')
				lineTotal += wordSpaces
				charSums = append(charSums, wordSpaces)
				cooked = append(cooked, '%')

				redraw()
				continue
			}

			// -----------------------------------------------------
			// LETTER
			// -----------------------------------------------------
			m, ok := MorseTable[r]
			if !ok {
				continue
			}

			// Inter-character gap → '#'
			if len(lineRunes) > 0 && lineRunes[len(lineRunes)-1] != ' ' {
				lineTotal += interCharSpaces
				charSums = append(charSums, interCharSpaces)
				cooked = append(cooked, '#')

				if lineTotal > MAX_WORD_SEGMENTS {
					lineTotal -= interCharSpaces
					charSums = charSums[:len(charSums)-1]
					cooked = cooked[:len(cooked)-1]
					fmt.Print("\r** TOO LONG — BACKSPACE TO FIX **")
					continue
				}
			}

			// Add the letter
			sum := m.Dark + m.Light
			lineTotal += sum
			charSums = append(charSums, sum)
			lineRunes = append(lineRunes, r)
			cooked = append(cooked, r)

			if lineTotal > MAX_WORD_SEGMENTS {
				lineTotal -= sum
				charSums = charSums[:len(charSums)-1]
				lineRunes = lineRunes[:len(lineRunes)-1]
				cooked = cooked[:len(cooked)-1]
				fmt.Print("\r** TOO LONG — BACKSPACE TO FIX **")
				continue
			}

			redraw()
		}
	}
}
