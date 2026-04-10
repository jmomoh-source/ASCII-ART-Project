package ascii

import "strings"

// VisibleLen calculates length of a string ignoring ANSI escape codes
func VisibleLen(s string) int {
	lenVisible := 0
	inEscape := false
	for _, r := range s {
		if r == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		lenVisible++
	}
	return lenVisible
}

// AlignLine applies alignment to a single line of text, accounting for ANSI codes
func AlignLine(line string, align string, width int) string {
	lineLen := VisibleLen(line)

	if lineLen >= width {
		return line
	}

	switch align {
	case AlignRight:
		return strings.Repeat(" ", width-lineLen) + line

	case AlignCenter:
		padding := (width - lineLen) / 2
		return strings.Repeat(" ", padding) + line

	default: // left
		return line
	}
}
