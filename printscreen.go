package main

// used to let user print instructions page or pattern map

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

type Screen struct {
	Lines []string
}

var reAnsi = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func openBrowser(path string) error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", path).Start()
	case "darwin":
		return exec.Command("open", path).Start()
	default:
		return exec.Command("xdg-open", path).Start()
	}
}

func (s *Screen) PrintScreen() error {
	var b strings.Builder

	b.WriteString(`<html><body style="margin:0; padding:0;">
    <pre style="margin:0; padding:0; font-family: monospace; font-size: 14px; white-space: pre-wrap; line-height: 1;">`)

	var lastWasBlank bool

	for _, line := range s.Lines {
		clean := strings.TrimRight(line, "\r\n")

		if visuallyBlank(clean) {
			if lastWasBlank {
				continue
			}
			lastWasBlank = true
		} else {
			lastWasBlank = false
		}

		htmlLine := ansiToHTML(clean)
		htmlLine = strings.TrimRight(htmlLine, "\r\n")

		b.WriteString(htmlLine)
		b.WriteString("\n")
	}

	b.WriteString(`</pre>
<script>
window.onload = function() {
    window.print();
};
window.onafterprint = function() {
    window.close();
};
</script>
</body></html>`)

	tmp := "printscreen.html"
	if err := os.WriteFile(tmp, []byte(b.String()), 0644); err != nil {
		return err
	}

	return openBrowser(tmp)
}

func visuallyBlank(s string) bool {
	// Remove ANSI sequences
	noAnsi := reAnsi.ReplaceAllString(s, "")
	// Trim whitespace
	noAnsi = strings.TrimSpace(noAnsi)
	return noAnsi == ""
}

// Drop-in replacements for fmt.Print, fmt.Println, fmt.Printf
// They print to the terminal AND capture the output in screen.Lines.
/*
func Println(a ...any) {
    text := fmt.Sprintln(a...)
    clean := strings.TrimRight(text, "\r\n")
    screen.Lines = append(screen.Lines, clean)
    fmt.Print(text) // Println behavior preserved
}

func Print(a ...any) {
    text := fmt.Sprint(a...)
    clean := strings.TrimRight(text, "\r\n")
    screen.Lines = append(screen.Lines, clean)
    fmt.Print(text)
}

func Printf(format string, a ...any) {
    text := fmt.Sprintf(format, a...)
    clean := strings.TrimRight(text, "\r\n")
    screen.Lines = append(screen.Lines, clean)
    fmt.Print(text)
}
*/

/*
func (s *Screen) Println(text string) {
    fmt.Println(text)
    s.Lines = append(s.Lines, text)
}
*/

func (s *Screen) Print(text string) {
	fmt.Print(text)
	s.Lines = append(s.Lines, text)
}

func (s *Screen) Printf(format string, a ...any) {
	text := fmt.Sprintf(format, a...)
	clean := strings.TrimRight(text, "\r\n")
	s.Lines = append(s.Lines, clean)
	fmt.Print(text)
}
