package main

import (
	"fmt"
	"os"
	"strings"

	"ascii-art/ascii"
)

func usage() {
	fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  go run . \"hello\"")
	fmt.Println("  go run . \"hello\" shadow")
	fmt.Println("  go run . --color=red \"hello\"")
	fmt.Println("  go run . --color=red \"He\" \"hello\"")
	fmt.Println("  go run . --output=out.txt \"hello\" standard")
	fmt.Println("  go run . --align=center \"hello\" standard")
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		usage()
		return
	}

	colorRules := make(map[string]string)
	var colorFlags []string

	outputFile := ""
	alignment := ""
	text := ""
	bannerType := ""

	// We iterate through arguments
	for i := 0; i < len(args); i++ {
		arg := args[i]

		if strings.HasPrefix(arg, "--color=") {
			color := strings.TrimPrefix(arg, "--color=")
			if color == "" {
				usage()
				return
			}
			colorFlags = append(colorFlags, color)
		} else if strings.HasPrefix(arg, "--output=") {
			outputFile = strings.TrimPrefix(arg, "--output=")
		} else if strings.HasPrefix(arg, "--align=") {
			alignment = strings.TrimPrefix(arg, "--align=")
		} else if strings.HasPrefix(arg, "--") {
			// Catch things like --color without =
			usage()
			return
		} else {
			// Positional
			if len(colorFlags) > 0 {
				hasMorePositional := false
				for j := i + 1; j < len(args); j++ {
					if !strings.HasPrefix(args[j], "--") {
						hasMorePositional = true
						break
					}
				}
				if hasMorePositional && text == "" {
					lastColor := colorFlags[len(colorFlags)-1]
					colorFlags = colorFlags[:len(colorFlags)-1]
					colorRules[arg] = lastColor
				} else if text == "" {
					text = arg
					for _, c := range colorFlags {
						colorRules[""] = c
					}
					colorFlags = nil
				} else {
					bannerType = arg
				}
			} else {
				if text == "" {
					text = arg
				} else if bannerType == "" {
					bannerType = arg
				}
			}
		}
	}

	for _, c := range colorFlags {
		colorRules[""] = c
	}

	if text == "" {
		usage()
		return
	}

	if alignment == "" {
		alignment = ascii.AlignLeft
	}

	width := ascii.GetTerminalWidth()
	
	result, err := ascii.GenerateAsciiArt(text, bannerType, colorRules, alignment, width)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if outputFile != "" {
		err := os.WriteFile(outputFile, []byte(result), 0644)
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
		}
	} else {
		// Output to terminal
		fmt.Print(result)
	}
}
