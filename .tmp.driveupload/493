package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// called PRIOR to getting a message so no idea of content
func getWordSpaces(modeFlg rune, hasSpaces bool, minWs, maxWs int) int {
	reader := bufio.NewReader(os.Stdin)

	for {
		prompt := ""

		if modeFlg == 'P' {
			prompt = fmt.Sprintf("Spaces IN BETWEEN words (from %d to %d).", minWs, maxWs)
		} else {
			// some morse mode E or M
			prompt = fmt.Sprintf("Spaces IN BETWEEN words and after the last one (from %d to %d, standard 7).", minWs, maxWs)
		}

		bar := strings.Repeat("=", len(prompt)+4)
		fmt.Println(bar)
		fmt.Printf("= %s =\n", prompt)
		fmt.Printf("%s\n > ", bar)
		ans, _ := reader.ReadString('\n')
		ans = strings.TrimSpace(ans)
		ws, err := strconv.Atoi(ans)
		if err != nil {
			fmt.Println("Invalid number:", ans)
			continue
		}
		if ws < minWs || ws > maxWs {
			fmt.Println("Invalid number:", ans)
			continue
		}
		return ws
	}
}
