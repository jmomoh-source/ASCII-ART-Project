package main

import (
	"fmt"
	"os"

	"ascii-art/ascii"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . [STRING]")
		return
	}

	input := os.Args[1]
	if input == "" {
		return
	}
	fmt.Print(ascii.AsciiArt(input))
}
