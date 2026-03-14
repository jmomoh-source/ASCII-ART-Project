package ascii

import (
	"fmt"
	"os"
	"strings"
)

// AsciiArt takes an input string and returns its ASCII art representation
// using the shadow, standard, banner file.
func AsciiArt(input string, banner string) string {
	if banner == "" {
		banner = "shadow"
	}
	inputFile, err := os.ReadFile("template/" + banner + ".txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return ""
	}

	content := strings.ReplaceAll(string(inputFile), "\r\n", "\n")
	inputFileLines := strings.Split(content, "\n")

	words := strings.Split(input, "\\n")
	result := ""

	for _, word := range words {
		if word == "" {
			result += "\n"
			continue
		}
		for i := 0; i < 8; i++ {
			for _, char := range word {
				result += inputFileLines[i+(int(char-' ')*9)+1]
			}
			result += "\n"
		}
	}
	return result
}
