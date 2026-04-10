package ascii

import "fmt"

// AsciiArt takes an input string and returns its ASCII art representation
// using the shadow, standard, banner file.
func AsciiArt(input string) string {
	res, err := GenerateAsciiArt(input, DefaultBanner, nil, AlignLeft, DefaultWidth)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return res
}
