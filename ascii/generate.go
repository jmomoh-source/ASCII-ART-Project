package ascii

import (
	"fmt"
	"strings"
)

// GenerateAsciiArt generates the full ASCII art string including banner selection, alignment, and coloring
func GenerateAsciiArt(input string, banner string, colorRules map[string]string, align string, width int) (string, error) {
	if align != "" && !IsValidAlignment(align) {
		return "", fmt.Errorf("invalid alignment setting")
	}
	if banner == "" {
		banner = DefaultBanner
	}

	for _, colorName := range colorRules {
		if _, ok := ParseColor(colorName); !ok {
			return "", fmt.Errorf("invalid color: %s", colorName)
		}
	}

	lines, err := LoadBanner(banner)
	if err != nil {
		return "", err
	}

	// Validate banner size
	if len(lines) < 855 {
		// Thinkertoy uses \r\n, LoadBanner splits by \n after replacing \r\n
	}

	words := strings.Split(input, "\\n")
	var result strings.Builder

	for _, word := range words {
		if word == "" {
			result.WriteString("\n")
			continue
		}

		// compute colors for this word
		charColors := GetCharColors(word, colorRules)
		
		var outputLines []string

		if align == AlignJustify {
			// Specific justify logic distributing spaces between words
			parts := strings.Split(word, " ")
			wordBlocks := make([][]string, len(parts))
			
			for w, part := range parts {
				wordBlocks[w] = renderWord(part, lines, GetCharColors(part, colorRules))
			}

			if len(parts) <= 1 {
				outputLines = wordBlocks[0]
			} else {
				// Width calculation
				wordWidths := make([]int, len(wordBlocks))
				totalWordWidth := 0
				for w, block := range wordBlocks {
					maxW := 0
					for _, line := range block {
						length := VisibleLen(line)
						if length > maxW {
							maxW = length
						}
					}
					wordWidths[w] = maxW
					totalWordWidth += maxW
				}

				totalSpaces := width - totalWordWidth
				gaps := len(parts) - 1
				if gaps > 0 {
					if totalSpaces < gaps {
						totalSpaces = gaps
					}
					evenSpace := totalSpaces / gaps
					extraSpace := totalSpaces % gaps

					outputLines = make([]string, CharHeight)
					for row := 0; row < CharHeight; row++ {
						lineResult := ""
						for w, block := range wordBlocks {
							rowContent := block[row]
							paddingNeeded := wordWidths[w] - VisibleLen(rowContent)
							padded := rowContent + strings.Repeat(" ", paddingNeeded)
							lineResult += padded

							if w < gaps {
								lineResult += strings.Repeat(" ", evenSpace)
								if w < extraSpace {
									lineResult += " "
								}
							}
						}
						outputLines[row] = lineResult
					}
				}
			}
		} else {
			outputLines = renderWord(word, lines, charColors)
			for i := range outputLines {
				outputLines[i] = AlignLine(outputLines[i], align, width)
			}
		}

		for _, line := range outputLines {
			result.WriteString(line)
			result.WriteString("\n")
		}
	}

	return result.String(), nil
}

func renderWord(word string, lines []string, charColors []string) []string {
	outputLines := make([]string, CharHeight)
	for i := 0; i < CharHeight; i++ {
		var rowBuilder strings.Builder
		for j, char := range word {
			if char < ' ' || char > '~' {
				continue
			}
			index := int(char-' ')*9 + 1
			if index+i < len(lines) {
				line := lines[index+i]
				if charColors != nil && j < len(charColors) && charColors[j] != "" {
					rowBuilder.WriteString(charColors[j])
					rowBuilder.WriteString(line)
					rowBuilder.WriteString(Colors["reset"])
				} else {
					rowBuilder.WriteString(line)
				}
			}
		}
		outputLines[i] = rowBuilder.String()
	}
	return outputLines
}
