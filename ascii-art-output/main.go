package main

import (
	"fmt"
	"os"
	"strings"

	"ascii-art/ascii"
)

func main() {
	args := os.Args[1:]
	var outputFile string

	if len(args) > 0 && strings.HasPrefix(args[0], "--output=") {
		outputFile = strings.TrimPrefix(args[0], "--output=")
		args = args[1:]
	}

	if len(args) < 1 || len(args) > 2 {
		fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]\n\nEX: go run . --output=<fileName.txt> something standard")
		return
	}

	input := args[0]
	if input == "" {
		return
	}

	banner := "standard"
	if len(args) == 2 {
		banner = args[1]
	}

	art := ascii.AsciiArt(input, banner)
	if art == "" {
		return
	}

	if outputFile != "" {
		err := os.WriteFile(outputFile, []byte(art), 0644)
		if err != nil {
			fmt.Println("Error writing file:", err)
		}
	} else {
		fmt.Print(art)
	}
}
