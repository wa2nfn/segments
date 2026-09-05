package main

import (
	"regexp"
	"strings"
)

// Converts ANSI-colored text to HTML
func ansiToHTML(outData string) string {
	re := regexp.MustCompile(`\x1b\[([0-9;]+)m`)
	var html strings.Builder
	var fgColor, bgColor string
	var openSpan bool

	for i := 0; i < len(outData); {
		if outData[i] == '\x1b' && i+1 < len(outData) && outData[i+1] == '[' {
			loc := re.FindStringIndex(outData[i:])
			if loc != nil {
				seq := outData[i : i+loc[1]]
				codes := re.FindStringSubmatch(seq)[1]
				for _, code := range splitCodes(codes) {
					switch {
					case code == "0":
						if openSpan {
							html.WriteString("</span>")
							openSpan = false
						}
						fgColor, bgColor = "", ""

					case strings.HasPrefix(code, "38:"):
						if openSpan {
							html.WriteString("</span>")
							openSpan = false
						}
						fgColor = parseRGB(code)

					case strings.HasPrefix(code, "48:"):
						if openSpan {
							html.WriteString("</span>")
							openSpan = false
						}
						bgColor = parseRGB(code)
					}
				}
				i += loc[1]
				continue
			}
		}

		ch := outData[i]
		escaped := htmlEscape(string(ch))

		if fgColor != "" || bgColor != "" {
			style := ""
			if fgColor != "" {
				style += "color:" + fgColor + ";"
			}
			if bgColor != "" {
				style += "background-color:" + bgColor + ";"
			}

			if fgColor == bgColor {
				html.WriteString(`<span style="` + style + `">█</span>`)
			} else {
				html.WriteString(`<span style="` + style + `">` + escaped + `</span>`)
			}
		} else {
			html.WriteString(escaped)
		}
		i++
	}

	return `<html><body style="margin:0; padding:0;"><pre style="margin:0; padding:0; font-family: monospace; font-size: 14px; white-space: pre-wrap; line-height: 1;">` + html.String() + `</pre></body></html>`
}
