package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/jwalton/go-supportscolor"
)

type Morse struct {
	Morse   string
	DLimage string
	Cimage  string
	Dark    int
	Light   int
}

type Result struct {
	S               int
	Angle           float64
	SEL             float64
	bdWidth         float64
	SwasFound       bool
	PaddingExterior float64
	PaddingInterior float64
}

// GroupStats holds counts of sequences by length for a rune
type GroupStats struct {
	Counts map[int]int // key = length (1,2,3), value = count
}

// Result holds overall analysis
type ResultCombo struct {
	SizeFlag  int                  // max group size flag (3, 2, or 1)
	Stats     map[rune]*GroupStats // per rune stats
	MinGroups int                  // total minimal groups needed
	Reduction int                  // total minimal groups needed
}

var stdWedges = map[int]struct{}{
	10: struct{}{},
	12: struct{}{},
	16: struct{}{},
	18: struct{}{},
	20: struct{}{},
	24: struct{}{},
	30: struct{}{},
	36: struct{}{},
	40: struct{}{},
}

// A struct to hold both the Foreground (FG) and Background (BG) RGB values
type TrueColor struct {
	R_FG, G_FG, B_FG int
	R_BG, G_BG, B_BG int
}

var screen Screen
var decimalPrecision = 32 // can override in wedges.txt
var kerf = 0.125          // can override in wedges.txt
var colorType string = "trueColor"
var paddingInterior float64
var paddingExterior float64
var SELmn string
var physicalOuterFlatRadius float64
var physicalOuterDiameter float64

// Define the color scheme using the TrueColor struct
var colorMap = map[rune]TrueColor{

	'D': {
		R_FG: 148, G_FG: 98, B_FG: 48,
		R_BG: 148, G_BG: 98, B_BG: 48,
	},

	'O': {
		R_FG: 181, G_FG: 140, B_FG: 83,
		R_BG: 181, G_BG: 140, B_BG: 83,
	},

	'L': {
		R_FG: 220, G_FG: 220, B_FG: 190,
		R_BG: 220, G_BG: 220, B_BG: 190,
	},
}

const (
	MORSE_INTERWORD           = 5
	MAX_WORD_SEGMENTS         = 240
	MIN_SEL           float64 = 0.25

	// just turn on foreground
	ANSI_GREEN = "\033[32m"
	ANSI_RED   = "\033[31m"
	ANSI_BROWN = "\033[38;5;94m\033[48;5;94m"
	ANSI_BEIGE = "\033[38;5;173m\033[48;5;173m"
	ANSI_TAN   = "\033[38;5;130m\033[48;5;130m"

	ANSI_RESET          = "\x1b[0m"
	MORSE_MAX_INTERWORD = 10
	DEG_TO_RADIANS      = math.Pi / 180
)

var reader = bufio.NewReader(os.Stdin)
var itsBasket = false

func main() {
	var haveOther bool = false
	var completeMorse = false
	var DLimage string
	var Cimage string
	var segs int
	var itsMorse bool = true
	var totalDark int
	var totalLight int
	var totalOther int
	var totalSegments int
	var interCharSpaces int = 3 // assumes morse
	var spaceSymbol string = "L"
	var word string
	var cookedWord string
	var wordSpaces int = 7 // std 7 morse code
	var diameter float64
	var diameter_lower_ring float64
	var colorLevel = "none"
	var ans string
	var hasSpaces bool
	var doEM bool
	var mode string = "'D' - Design"

	if supportscolor.Stderr().Has16m {
		colorLevel = "trueColor"
	} else if supportscolor.Stdout().Has256 {
		colorLevel = "256Color"
	} else if supportscolor.Stdout().SupportsColor {
		colorLevel = "16Color"
	}

	ClearScreen()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(0)
	}()

	MorseTable := buildMorseTable()
	fmt.Print("\n\n\t\t\t                 Segments v1.0")
	fmt.Print("\n\n\t\t\tDesigning A Decorative Segmented Ring For A Bowl")
	fmt.Print("\n\t\t\t                         By")
	fmt.Print("\n\t\t\t                Bill Lanahan, WA2NFN\n\n")

	if colorLevel == "none" {
		fmt.Print("\nSorry this app requires appropriate color level support.\n")
		os.Exit(1)
	}

	fmt.Print("( Ctrl-C to quit. )\n")
	readWedgesFile()

	prompt := "Enter 'Y' for the instructions section."
	bar := strings.Repeat("=", len(prompt)+4)
	fmt.Println(bar)
	fmt.Printf("= %s =\n", prompt)
	fmt.Printf("%s\n > ", bar)

	reader := bufio.NewReader(os.Stdin)
	ins, _ := reader.ReadString('\n')
	ans = strings.TrimSpace(ins)

	if len(ans) > 0 && (ans[0] == 'y' || ans[0] == 'Y') {
		morseHelp()
	}

	prompt = "Enter 'Y' to see pattern details, RETURN to continue."
	bar = strings.Repeat("=", len(prompt)+4)
	fmt.Println(bar)
	fmt.Printf("= %s =\n", prompt)
	fmt.Printf("%s\n > ", bar)

	ans, _ = reader.ReadString('\n')
	ans = strings.TrimSpace(ans)
	ans = strings.ToUpper(ans)

	if len(ans) > 0 && ans[0] == 'Y' {
		showMorse(MorseTable, wordSpaces)

		var input string
		fmt.Print("\nNote: Default length of SPACE ' ' is shown.\n'*' represents the ProSign <SOS>.\ing/
		nter 'P' to print, RETURN to continue.")
		fmt.Scanln(&input)

		if strings.ToUpper(input) == "P" {
			fmt.Println()
			screen.PrintScreen()
		}
	}

	prompt = "Enter 'M' for MorseCode message, 'E' for Enhanced MorseCode, 'D' or RETURN for Design."
	bar = strings.Repeat("=", len(prompt)+4)
	fmt.Println(bar)
	fmt.Printf("= %s =\n", prompt)
	fmt.Printf("%s\n > ", bar)
	ans, _ = reader.ReadString('\n')
	ans = strings.TrimSpace(ans)
	ans = strings.ToUpper(ans)

	if len(ans) > 0 && ans[0] == 'M' {
		mode = "'M' - MorseCode"
		itsMorse = true
		completeMorse = true
	} else if len(ans) > 0 && ans[0] == 'E' {
		mode = "'E' - Enhanced MorseCode"
		itsMorse = true
		doEM = true
	} else {
		itsMorse = false
	}

	prompt = "Enter 'B' for Basket Illusion vs. Wood Segments, RETURN for wood segments."
	bar = strings.Repeat("=", len(prompt)+4)
	fmt.Println(bar)
	fmt.Printf("= %s =\n", prompt)
	fmt.Printf("%s\n > ", bar)
	ans, _ = reader.ReadString('\n')
	ans = strings.TrimSpace(ans)
	ans = strings.ToUpper(ans)
	if len(ans) > 0 && ans[0] == 'B' {
		mode = "'B' - Basket Illusions"
		itsBasket = true
	}

	// get space requirements PRIOR to getting message

	if !itsMorse {
		for {
			prompt = "Number of inter-pattern(character) spaces (1-6)"
			bar = strings.Repeat("=", len(prompt)+4)
			fmt.Println(bar)
			fmt.Printf("= %s =\n", prompt)
			fmt.Printf("%s\n > ", bar)
			ans, _ = reader.ReadString('\n')
			ans = strings.TrimSpace(ans)
			input, err := strconv.Atoi(ans)

			if err != nil {
				fmt.Println("Invalid number: ", ans)
				continue
			}
			if input < 1 || input > 6 {
				fmt.Printf("\n*** ERROR: number of spaces is invalid.\n")
				continue
			} else {
				interCharSpaces = input
				break
			}
		}

		wordSpaces = getWordSpaces('D', hasSpaces, interCharSpaces+1, MORSE_MAX_INTERWORD) // let them has a little as 1

		prompt = "Enter 'Y' for OTHER wood for the spaces IN BETWEEN characters, else RETURN."
		bar = strings.Repeat("=", len(prompt)+4)
		fmt.Println(bar)
		fmt.Printf("= %s =\n", prompt)
		fmt.Printf("%s\n > ", bar)
		ans, _ = reader.ReadString('\n')
		ans = strings.TrimSpace(ans)
		ans = strings.ToUpper(ans)

		if len(ans) > 0 && ans[0] == 'Y' {
			haveOther = true
			spaceSymbol = "O"
		}
	} else {
		// end space is required
		// its E or M
		wordSpaces = getWordSpaces('M', hasSpaces, 4, 40) // made 4 very close to char space. Visually ??
	}

	// update the Map for space char
	SpaceTable := MorseTable[' ']

	tmp := 0
	if wordSpaces != 0 {
		tmp = wordSpaces
		SpaceTable.Morse = strings.Repeat(" ", tmp)
	} else {
		tmp = interCharSpaces
		SpaceTable.Morse = strings.Repeat(" ", tmp)
	}

	if haveOther {
		SpaceTable.DLimage = strings.Repeat("O", tmp)
	} else {
		SpaceTable.DLimage = strings.Repeat("L", tmp)
	}
	SpaceTable.Light = tmp

	SpaceTable.Dark = 0
	MorseTable[' '] = SpaceTable

	// always collect regardless of mode
	var err error
	for {
		adj := "message"
		if !itsMorse {
			adj = "pattern details"
		}
		prompt = fmt.Sprintf("Enter \"%s\" (allowed: A-Z0-9 .,/?=+-:;'\"@*) (max. %d segments., total in <> below).", adj, MAX_WORD_SEGMENTS)
		bar = strings.Repeat("=", len(prompt)+4)
		fmt.Println(bar)
		fmt.Printf("= %s =\n", prompt)
		fmt.Printf("%s\n > ", bar)

		// readMessage gets "word" by reading keyboard
		word, cookedWord, err = readMessageCooked(MorseTable, wordSpaces, interCharSpaces, itsMorse)

		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(22)
		}

		if len(word) > MAX_WORD_SEGMENTS {
			fmt.Println("\n** Error: message word too long.")
			os.Exit(1)
		} else if len(word) < 1 {
			fmt.Println("\n** Error: message word empty.")
		} else {
			break
		}
	}

	// now all spacing is know from user and message/pattern
	// adds slight flexibility to MAYBE minimize segments
	if strings.Contains(word, " ") {
		hasSpaces = true
	}

	interCharImageDL := strings.Repeat(spaceSymbol, interCharSpaces)
	interCharImageC := strings.Repeat(spaceSymbol, interCharSpaces)

	ClearScreen()
	fmt.Printf("\n\n Parsing Message/Design By Character\n ===================================\n")

	// we have read the message from the user - too late for his corrections if needed
	// now we reparse for the purpose of calculationg and adjusting the end
	var changeCandidates []int
	runes := []rune(cookedWord)
	length := len(runes)

	fmt.Printf(" Name   Segments       Pattern as D&L           Colorized\n")
	fmt.Printf(" ----   --------   ----------------------   -----------------------\n")
	for i := 0; i < length; i++ {
		r := runes[i]
		val := MorseTable[r]

		// for visables
		if r != '#' && r != '%' { // plain rune
			val.Cimage = colorize(val.DLimage)
			fmt.Printf(" '%c'       %2d      %-23s  %s\n",
				r, val.Dark+val.Light, val.DLimage, val.Cimage)

			totalDark += val.Dark
			totalLight += val.Light
			DLimage += val.DLimage
		} else if r == '#' { // space
			if haveOther {
				totalOther += interCharSpaces
				fmt.Printf("%-10s %2d\n", "Other", interCharSpaces)
			} else {
				totalLight += interCharSpaces
				fmt.Printf("%-10s %2d\n", "space", interCharSpaces)
			}
			DLimage += interCharImageDL
			Cimage += interCharImageC
		} else if r == '%' { // space char
			nr := ' '
			nval := MorseTable[nr]

			if !itsMorse {
				spaceColorImage := colorize(nval.DLimage)
				if haveOther {
					totalOther += nval.Light
				} else {
					totalLight += nval.Light //wordSpaces
				}
				fmt.Printf(" '%c'       %2d      %-23s  %s\n",
					nr, nval.Light, nval.DLimage, spaceColorImage)
				DLimage += nval.DLimage
				Cimage += spaceColorImage
			} else {
				// canonical wordspace
				totalLight += nval.Light
				DLimage += nval.DLimage
				nval.Cimage = colorize(nval.DLimage)
				Cimage += nval.Cimage
				fmt.Printf(" '%c'       %2d      %-23s  %s\n",
					nr, nval.Dark+nval.Light, nval.DLimage, nval.Cimage)
			}
		}
	}

	// cookdWord is done, now append ending space
	if itsMorse {
		DLimage += strings.Repeat("L", wordSpaces)
		totalLight += wordSpaces
		fmt.Printf("%-10s %2d\n", "end space", wordSpaces)
	} else {
		endSpaces := interCharSpaces // its a pattern
		if haveOther {
			DLimage += strings.Repeat("O", endSpaces)
			totalOther += endSpaces
		} else {
			DLimage += strings.Repeat("L", endSpaces)
			totalLight += endSpaces
		}
		fmt.Printf("%-10s %2d\n", "end space", endSpaces)
	}

	totalSegments = totalDark + totalLight + totalOther
	fmt.Printf("\n           %2d Total before adjustment\n", totalSegments)

	//
	// special case to match segment count
	//
	if completeMorse && !itsBasket {
		fmt.Println()
		for {
			answer := ""
			prompt = "ONLY if you are matching another ring's SEL., enter the number of EXTRA =\n=segments to add to this ring, RETURN to continue: "
			bar := strings.Repeat("=", 73)
			fmt.Println(bar)
			fmt.Printf("= %s \n", prompt)

			// Read user input BEFORE doing anything else
			reader := bufio.NewReader(os.Stdin)
			answer, _ = reader.ReadString('\n')
			answer = strings.TrimSpace(answer)

			if answer == "" {
				break
			}

			num, err := strconv.Atoi(answer)
			if err != nil || num < 0 {
				fmt.Print("\n *** Error: invalid input. 0 or a positive number.\n")
				continue
			}

			if num == 0 {
				break
			}

			if num+totalSegments > MAX_WORD_SEGMENTS {
				fmt.Print("\n *** Error: invalid input. TotalSegments exceed maximum.\n")
				continue
			}

			// add to segment count and the DLimage
			DLimage += strings.Repeat("L", num)
			totalSegments += num
			totalLight += num
			break
		}
	}

	// split for Basket
	if !itsBasket {
		if totalSegments < 3 {
			fmt.Printf("\n*** Error: Segments must be between 3 and %d. Your's is %d.\n", MAX_WORD_SEGMENTS, totalSegments)

			if itsMorse && !doEM {
				fmt.Print("\nTry different input pattern.\n\n")
			}
			os.Exit(2)
		} else if totalSegments > 2*180 {
			fmt.Printf("\n*** Error: Segments must be between 3 and 180. Your's is %d.\n", totalSegments)
			fmt.Print("\nA segment count exceeding 180 makes a cut angle of less than 1 degree.")

			if itsMorse && !doEM {
				fmt.Print("\nTry different input message or use the Enhanced MorseCode mode.\n\n")
			} else {
				fmt.Print("\nTry different input pattern.\n\n")
			}
			os.Exit(2)
		}
	}

	// move on if NOT basket illusions
	// get user inputs and check them
	var ringWidth float64
	var result Result

	fmt.Printf("\n\n Geometry Inputs\n ===============\n")
	diameter = validateInput("Final ring DIAMETER ", 4.0, 24.0)
	if !itsBasket {
		diameter_lower_ring = validateInput("Ring diameter LOWER RING", 2.0, 24.0)
		ringWidth = validateInput("Ring WIDTH (wall thickness) ", 0.25, 3.0)
		paddingExterior = validateInput("Ring PADDING EXTERIOR (ShopSlop)", 0, 1.0)
		paddingInterior = validateInput("Ring PADDING INTERIOR (ShopSlop)", 0, 1.0)
	} else {
		diameter_lower_ring = diameter
		ringWidth = 5.0
		paddingExterior = 0
		paddingInterior = 0
	}

	// segmented bowl
	result, err = CalculatePolygonRing(diameter, diameter_lower_ring, ringWidth, paddingExterior, paddingInterior, totalSegments)

	if err != nil {
		if !itsBasket {
			fmt.Printf("*** Error: inputs diameter=%.3f, ringWidth=%.3f, paddingExterior=%0.3f, paddingInterior=%0.3f, segments=%d\n\n", diameter, ringWidth, paddingExterior, paddingInterior, result.S)
		} else {
			fmt.Printf("*** Error: input diameter=%.3f\n\n", diameter)
		}
		os.Exit(66)
	}

	segs = result.S

	var addedSpaces = segs - totalSegments

	diff := 0
	if totalSegments < segs {
		diff = segs - totalSegments
		totalLight += diff
		if haveOther {
			DLimage += strings.Repeat("O", diff)
		} else {
			DLimage += strings.Repeat("L", diff)
		}
	} else if totalSegments > segs {
		diff = totalSegments - segs
		totalLight -= diff
		DLimage = DLimage[:len(DLimage)-diff]
	}

	var delta int
	if doEM {
		DLimage, delta = adjustLs(DLimage, wordSpaces)
	}

	if !itsMorse {
		DLimage, delta = pModeAdj(DLimage, interCharSpaces, changeCandidates)
	}

	// display required details in an obvious manner
	var outData string
	SELmn = toMixedNumber(result.SEL)

	// danger test
	SELcheck(result.SEL)

	SegWidth := toMixedNumber(result.bdWidth)
	var boxLine1, boxLine2, boxLine3 string

	if haveOther && totalOther > 0 {
		boxLine1 = fmt.Sprintf("Total Segs: %d%s                  DARK Segs: %d, LIGHT Segs: %d, Other Segs: %d", segs, mkStd(segs), totalDark, totalLight, totalOther)
	} else {
		boxLine1 = fmt.Sprintf("Total Segs: %d%s                  DARK Segs: %d, LIGHT Segs: %d.", segs, mkStd(segs), totalDark, totalLight)
	}

	if itsBasket {
		boxLine2 = ""
		boxLine3 = fmt.Sprintf("SEL <%s>", toMixedNumber(result.SEL))
	} else {
		boxLine2 = fmt.Sprintf("Angle: Seg. Miter Angle: %0.2f deg., WEDGIE Central Angle: %0.2f deg.", result.Angle, result.Angle*2.0)
		boxLine3 = fmt.Sprintf("SEL <%s>, SegWidth <%s>", toMixedNumber(result.SEL), SegWidth)
	}

	outData += BuildSummaryBox(" RING SUMMARY ", boxLine1, boxLine2, boxLine3)
	if itsMorse {
		outData += fmt.Sprintf("\n\n Inputs:\n\tMode: <%s>\n\tWord Spaces: <%d>", mode, wordSpaces)
	} else {
		if wordSpaces == 0 {
			outData += fmt.Sprintf("\n\n Inputs:\n\tMode: <%s>\n\tInter Char Spaces: <%d>", mode, interCharSpaces)
		} else {
			outData += fmt.Sprintf("\n\n Inputs:\n\tMode: <%s>\n\tInter Char Spaces: <%d>, Word Spaces: <%d>", mode, interCharSpaces, wordSpaces)
		}
	}

	if itsBasket {
		outData += fmt.Sprintf("\n\tDiameter: <%s>\n", toMixedNumber(diameter))
	} else {
		outData += fmt.Sprintf("\n\tDiameter: <%s>, Ring Width: <%s>\n\tPadding Exterior: <%s>, Padding Interior: <%s>\n", toMixedNumber(diameter), toMixedNumber(ringWidth), toMixedNumber(paddingExterior), toMixedNumber(paddingInterior))
	}

	outData += fmt.Sprintf("\n Pattern/message: <%s> (plus end spaces)\n", word)
	if addedSpaces > 0 {
		if delta > 1 {
			outData += fmt.Sprintf("\n Note: Segment diagram shows (%d) added spaces for angle adjustment.\n <%d> not redistributed.\n", addedSpaces, delta)
		} else if delta < 0 {
			outData += fmt.Sprintf("\n Note: Segment diagram shows (%d) added spaces for angle adjustment.\n Insufficient qty. for redistribution.\n", addedSpaces)
		}
	}

	str, rc := printNumLine(DLimage)
	outData += "\n Below are idices of the segments above the ring in both Dark and Light segments\n as well as in colorized format. Use these to glue up your ring.\n\n"

	outData += str

	if rc {
		outData += "IMPORTANT: The second set of lines directly extend the first set. Do not leave a gap or forget\nthe second set if you segments count exceeded 120!\n\n"
	}
	bdLengthDark := (float64(totalDark) * float64(result.SEL)) + (kerf * float64(totalDark))
	bdLengthLight := (float64(totalLight) * float64(result.SEL)) + (kerf * float64(totalLight))

	if !itsBasket {
		outData += "\n\n Board Material:\n"

		Dtile := colorize("D")
		Ltile := colorize("L")
		outData += fmt.Sprintf(" %s %s Dark Board length <%.2f> inches approximately.\n", Dtile, ANSI_RESET, math.Ceil(bdLengthDark))

		outData += fmt.Sprintf(" %s %s Light Board length <%.2f> inches approximately.\n", Ltile, ANSI_RESET, math.Ceil(bdLengthLight))

		if haveOther && totalOther > 0 {
			bdLengthOther := (float64(totalOther) * float64(result.SEL)) + (kerf * float64(totalOther))
			Otile := colorize("O")
			outData += fmt.Sprintf(" %s %s Other Board length <%.2f> inches approximately.\n", Otile, ANSI_RESET, math.Ceil(bdLengthOther))
		}

		outData += fmt.Sprintf("\n          board width <%s> inches (includes padding).\n", toMixedNumber(result.bdWidth))

		outData += fmt.Sprintf("\n          * assumes %v\" kerf\n", toMixedNumber(kerf))
		outData += "          ** add enough to cut safely (approx 2+ inches)."

		important(outData)
	}
	fmt.Printf("%s%s\n\n", outData, ANSI_RESET)

	prompt = "Enter 'Y' to view REVERSAL of Dark and Light segments, RETURN to skip."
	bar = strings.Repeat("=", len(prompt)+4)
	fmt.Println(bar)
	fmt.Printf("= %s =\n", prompt)
	fmt.Printf("%s\n > ", bar)
	ans, _ = reader.ReadString('\n')
	ans = strings.TrimSpace(ans)

	if len(ans) > 0 && (ans[0] == 'y' || ans[0] == 'Y') {

		// display required details in an obvious manner
		var outDataRev string
		SELmn = toMixedNumber(result.SEL)

		// danger test
		SELcheck(result.SEL)

		SegWidth = toMixedNumber(result.bdWidth)

		if haveOther && totalOther > 0 {
			boxLine1 = fmt.Sprintf("Total Segs: %d%s                  DARK Segs: %d, LIGHT Segs: %d, Other Segs: %d", segs, mkStd(segs), totalDark, totalLight, totalOther)
		} else {
			boxLine1 = fmt.Sprintf("Total Segs: %d%s                  DARK Segs: %d, LIGHT Segs: %d.", segs, mkStd(segs), totalDark, totalLight)
		}

		if !itsBasket {
			boxLine2 = fmt.Sprintf("Angle: Seg. Miter Angle: %0.2f deg., WEDGIE Central Angle: %0.2f deg.", result.Angle, result.Angle*2.0)
			boxLine3 = fmt.Sprintf("SEL <%s>, SegWidth <%s>", toMixedNumber(result.SEL), SegWidth)
		} else {
			boxLine2 = "                                                                   "
			boxLine3 = fmt.Sprintf("Width <%s>", toMixedNumber(result.SEL))
		}

		outDataRev += BuildSummaryBox(" RING SUMMARY ", boxLine1, boxLine2, boxLine3)
		if itsMorse {
			outDataRev += fmt.Sprintf("\n\n Inputs:\n\tMode: <%s>\n\tWord Spaces: <%d>", mode, wordSpaces)
		} else {
			if wordSpaces == 0 {
				outDataRev += fmt.Sprintf("\n\n Inputs:\n\tMode: <%s>\n\tInter Char Spaces: <%d>", mode, interCharSpaces)
			} else {
				outDataRev += fmt.Sprintf("\n\n Inputs:\n\tMode: <%s>\n\tInter Char Spaces: <%d>, Word Spaces: <%d>", mode, interCharSpaces, wordSpaces)
			}
		}

		if !itsBasket {
			outDataRev += fmt.Sprintf("\n\tDiameter: <%s>, Ring Width: <%s>\n\tPadding Exterior: <%s>, Padding Interior: <%s>\n", toMixedNumber(diameter), toMixedNumber(ringWidth), toMixedNumber(paddingExterior), toMixedNumber(paddingInterior))
		} else {
			outDataRev += fmt.Sprintf("\n\tDiameter: <%s>\n", toMixedNumber(diameter))
		}

		outDataRev += fmt.Sprintf("\n Pattern/message: <%s> (plus end spaces)\n", word)
		if addedSpaces > 0 {
			if delta > 1 {
				outDataRev += fmt.Sprintf("\n Note: Segment diagram shows (%d) added spaces for angle adjustment.\n <%d> not redistributed.\n", addedSpaces, delta)
			} else if delta < 0 {
				outDataRev += fmt.Sprintf("\n Note: Segment diagram shows (%d) added spaces for angle adjustment.\n Insufficient qty. for redistribution.\n", addedSpaces)
			}
		}

		//
		//OrigDLimage = DLimage
		DLimage = swapChars(DLimage, 'D', 'L')
		var tL = totalLight
		var tD = totalDark
		totalDark = tL
		totalLight = tD

		str, rc = printNumLine(DLimage)
		outDataRev += "\n Below are idices of the segments above the ring in both Dark and Light segments\n as well as in colorized format. Use these to glue up your ring.\n\n"

		outDataRev += str

		if rc {
			outDataRev += "\n***IMPORTANT: The second set of lines directly extend the first set. Do not leave a gap or forget\nthe second set!\n\n"
		}

		if !itsBasket {
			bdLengthDark = (float64(totalDark) * float64(result.SEL)) + (kerf * float64(totalDark))
			bdLengthLight = (float64(totalLight) * float64(result.SEL)) + (kerf * float64(totalLight))

			outDataRev += "\n\n Board Material:\n\n"

			Dtile := colorize("D")
			Ltile := colorize("L")
			outDataRev += fmt.Sprintf(" %s %s  Dark Board length <%.2f> inches approximately.\n", Dtile, ANSI_RESET, math.Ceil(bdLengthDark))

			outDataRev += fmt.Sprintf(" %s %s Light Board length <%.2f> inches approximately.\n", Ltile, ANSI_RESET, math.Ceil(bdLengthLight))

			if haveOther && totalOther > 0 {
				bdLengthOther := (float64(totalOther) * float64(result.SEL)) + (kerf * float64(totalOther))
				Otile := colorize("O")
				outData += fmt.Sprintf(" %s %s Other Board length <%.2f> inches approximately.\n", Otile, ANSI_RESET, math.Ceil(bdLengthOther))
			}
			outDataRev += fmt.Sprintf("\n          board width <%s> inches (includes padding).\n", toMixedNumber(result.bdWidth))

			outDataRev += fmt.Sprintf("\n          * assumes %v\" kerf\n", toMixedNumber(kerf))
			outDataRev += "          ** add enough to cut safely (approx 2+ inches)"
		}

		fmt.Printf("%s%s\n\n", outDataRev, ANSI_RESET)

	}
	fmt.Print("\n")

	prompt = "File name to save to (NO suffix), (RETURN to exit): "
	bar = strings.Repeat("=", len(prompt)+4)
	reader = bufio.NewReader(os.Stdin)

	for {
		fmt.Println(bar)
		fmt.Printf("= %s =\n", prompt)
		fmt.Printf("%s\n > ", bar)

		ans, _ = reader.ReadString('\n')
		ans = strings.TrimSpace(ans)

		if ans != "" {
			break
		}

		// Strip any suffix the user typed
		base := filepath.Base(ans)
		ext := filepath.Ext(base)
		ans = strings.TrimSuffix(base, ext)

		// Add .html suffix
		ans = ans + ".html"
		// Check if file exists
		if _, err := os.Stat(ans); err == nil {
			// File exists — warn and ask for confirmation
			fmt.Printf("WARNING: file <%s> exists. Enter 'Y' to overwrite: ", ans)
			reply, _ := reader.ReadString('\n')
			reply = strings.TrimSpace(reply)

			if strings.ToUpper(reply) != "Y" {
				fmt.Println("Not overwriting. Please enter a different name.\n")
				continue // reprompt for filename
			}
		}
		break
	}

	// Convert ANSI to HTML
	outData = ansiToHTML(outData)

	// Attempt to write file
	err = os.WriteFile(ans, []byte(outData), 0644)
	if err != nil {
		fmt.Printf("Error: Failed to write file: %v\n", err)
		fmt.Print("Take a screen shot before closing the app.\n\n")
		os.Exit(1)
	}

	fmt.Printf("\nResults saved in file: %s\n\n", ans)

	prompt = fmt.Sprintf("Enter 'P' to print file <%s>, RETURN to continue: ", ans)
	bar = strings.Repeat("=", len(prompt)+4)
	fmt.Println(bar)
	fmt.Printf("= %s =\n", prompt)
	fmt.Printf("%s\n > ", bar)

	// Read user input BEFORE doing anything else
	reader = bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(answer)

	if strings.ToUpper(answer) == "P" {
		fmt.Print("\nAfter PRINT, you may need to close a blank browser page.\nPreparing file, standby ...")
		// Only now do we read the HTML and inject auto-print
		// read saved HTML
		htmlBytes, err := os.ReadFile(ans)
		if err != nil {
			fmt.Println("Error reading saved HTML:", err)
			os.Exit(1)
		}
		htmlContent := string(htmlBytes)

		// inject auto-print script
		autoPrint := `
<script>
window.onload = function() {
    window.print();
};
window.onafterprint = function() {
    document.body.innerHTML = "";
    document.title = "";
};
</script>
</body>`

		modified := strings.Replace(htmlContent, "</body>", autoPrint, 1)

		// write auto-print file
		autoFile, err := os.CreateTemp(os.TempDir(), "autoprint-*.html")
		if err != nil {
			fmt.Println("Error creating auto-print file:", err)
			os.Exit(1)
		}
		defer os.Remove(autoFile.Name())

		if _, err = autoFile.Write([]byte(modified)); err != nil {
			fmt.Println("Error writing auto-print file:", err)
			os.Exit(1)
		}

		autoFile.Close()
		openBrowser(autoFile.Name())
	}
	os.Exit(0)
}

func parseFractionalFloat(s string) (float64, error) {
	s = strings.TrimSpace(s)

	// Normalize hyphens: "12-1/2" → "12 1/2"
	s = strings.ReplaceAll(s, "-", " ")

	// Normalize spaces around slash
	// "1 / 2" → "1/2"
	// "1/ 2" → "1/2"
	// "1 /2" → "1/2"
	for strings.Contains(s, " /") || strings.Contains(s, "/ ") {
		s = strings.ReplaceAll(s, " /", "/")
		s = strings.ReplaceAll(s, "/ ", "/")
	}

	// If it contains a slash, treat it as a fraction or mixed fraction
	if strings.Contains(s, "/") {
		parts := strings.Fields(s)

		// Case 1: pure fraction "1/2"
		if len(parts) == 1 {
			numDen := strings.Split(parts[0], "/")
			if len(numDen) != 2 {
				return 0, fmt.Errorf("invalid fraction: %s", parts[0])
			}
			num, err1 := strconv.Atoi(numDen[0])
			den, err2 := strconv.Atoi(numDen[1])
			if err1 != nil || err2 != nil || den == 0 {
				return 0, fmt.Errorf("invalid fraction: %s", parts[0])
			}
			return float64(num) / float64(den), nil
		}

		// Case 2: mixed fraction "12 1/2"
		if len(parts) == 2 {
			wholeStr := parts[0]
			fracStr := parts[1]

			numDen := strings.Split(fracStr, "/")
			if len(numDen) != 2 {
				return 0, fmt.Errorf("invalid fraction: %s", fracStr)
			}

			whole, err0 := strconv.Atoi(wholeStr)
			num, err1 := strconv.Atoi(numDen[0])
			den, err2 := strconv.Atoi(numDen[1])
			if err0 != nil || err1 != nil || err2 != nil || den == 0 {
				return 0, fmt.Errorf("invalid mixed fraction: %s", s)
			}

			return float64(whole) + float64(num)/float64(den), nil
		}

		return 0, fmt.Errorf("invalid fractional format: %s", s)
	}

	// No slash normal float
	return strconv.ParseFloat(s, 64)
}

func SELcheck(sel float64) {
	if sel < 0.1 {
		fmt.Printf("*** ERROR: For your bowl diameter, the number of segments is\ncausing your SEL length to be less than 1/4 inch.\n")
		os.Exit(88)
	}
}

// final screen infor about dirsect of assembly
func important(outData string) {

	fmt.Print(`
*** IMPORTANT ***

The below ring segments are numbered from 1 to the ring length, LEFT 
to RIGHT. If the ring is a Morse Code message, then the correct DIRECTION
is critical, or the message will be a mirror image.

OUTSIDE VIEW: If the ring's position on the bowl, and the slope of the 
bowl's side at that position make it visable when viewed from the OUTSIDE,
then build the ring on your desktop **COUNTERCLOCKWISE** as shown below.

(Assuming your message word was SAMPLE, of course each letter represents all 
the segments for that letter. S would really be 5 segments plus spacing).
Lastly, since this is simplified explanation, there would also be the 
end of message space segments between 6 and 1.

            (6) E    L  (5)
	\  (1) S        P  (4) 
	 \ (2)  A   M  (3)
	  \_____> reading direction


INSIDE VIEW: If the ring's position on the bowl, and the slope of the 
bowl's side is shallow (more like a platter of soup bowl) then the 
message is better viewed from the INSIDE, then build the ring on your 
desktop **CLOCKWISE** as shown below.

	   _____> reading direction
          /
         /   (2) A    M  (3)
	/  (1) S        P  (4) 
	    (6)  E   L  (5)

`)

}

func validateInput(prompt string, min, max float64) float64 {
	for {
		fmt.Printf("%25s (min %.2f, max %.2f): ", prompt, min, max)

		ans, _ := reader.ReadString('\n')
		ans = strings.TrimSpace(ans)
		if len(ans) == 0 {
			fmt.Print("** Error: A value is required.\n")
			continue
		}

		val, err := parseFractionalFloat(ans)
		if err != nil {
			fmt.Println("** Error:", err)
			continue
		}
		if val < min || val > max {
			fmt.Printf("** Error: value must be between %.2f and %.2f\n", min, max)
			continue
		}
		return val
	}
}

func CalculatePolygonRing(diameter, under_ring_diameter, ringWidth, PaddingExterior float64, PaddingInterior float64, S int) (Result, error) {

	// width calculation

	Angle := 180.0 / float64(S)

	var found bool
	var SEL float64
	var BdWidth float64
	if !itsBasket {
		_, found = stdWedges[S]

		if !found {

			for s := S; s <= 180; s++ {
				angle := 180.0 / float64(s) // degrees
				nearestInt := math.Round(angle)
				diff := math.Abs(angle - nearestInt)

				if diff <= 0.1 {

					// Peek exactly ONE step ahead to see if it's a "perfect" fit (e.g., 59 vs 60)
					nextAngle := 180.0 / float64(s+1)
					nextDiff := math.Abs(nextAngle - math.Round(nextAngle))

					if nextDiff < diff {
						// The next segment count is a better fit, use it instead!
						S = s + 1
						Angle = math.Round(nextAngle)
					} else {
						// The current segment count is better, stick with it!
						S = s
						Angle = nearestInt
					}

					break
				}
			}
		}

		maxRadius := max(diameter/2.0, under_ring_diameter/2.0)
		minRadius := min((diameter/2.0 - ringWidth), (under_ring_diameter/2.0 - ringWidth))

		targetInnerRadius := minRadius - paddingInterior
		if targetInnerRadius <= 0 {
			return Result{}, errors.New("RingWidth and All Padding are too large for this diameter (results in negative or zero inner radius)")
		}

		angleRad := Angle * (math.Pi / 180.0)
		BdWidth = CalculateSegmentWidth(maxRadius, minRadius, paddingExterior, paddingInterior, S)

		physicalOuterFlatRadius = targetInnerRadius + BdWidth
		physicalOuterDiameter = physicalOuterFlatRadius * 2.0
		// Tangent flat-to-flat measurement
		SEL = physicalOuterDiameter * math.Tan(angleRad)

	} else {
		// itsBasket just need SEL
		SEL = (math.Pi * diameter) / float64(S)
	}

	return Result{
		S:               S,
		Angle:           Angle,
		SEL:             SEL,
		bdWidth:         BdWidth,
		SwasFound:       found,
		PaddingExterior: PaddingExterior,
		PaddingInterior: PaddingInterior,
	}, nil
}

func swapChars(s string, r1, r2 rune) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case r1:
			return r2
		case r2:
			return r1
		default:
			return r
		}
	}, s)
}

// for a segmented ring, properly accounting for ring slope, surface padding,
// and the inner polygon sagitta.
func CalculateSegmentWidth(rOutMax, rInMin, paddingExterior, paddingInterior float64, segs int) float64 {

	// 1. The outer flat edge (apothem) must sit at the max outer radius + padding
	outerApothem := rOutMax + paddingExterior

	// 2. The INNER CORNERS must sit at the min inner radius - padding.
	// We calculate the inner flat edge (apothem) by factoring in the polygon angle.
	angle := math.Pi / float64(segs)
	innerApothem := (rInMin - paddingInterior) * math.Cos(angle)

	// 3. The required board width is the distance between the two parallel cuts
	return outerApothem - innerApothem
}
