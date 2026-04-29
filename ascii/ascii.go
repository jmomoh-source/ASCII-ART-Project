package ascii

import "fmt"

func AsciiArt(input string) string {
	res, err := GenerateAsciiArt(input, DefaultBanner, nil, AlignLeft, DefaultWidth)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	return res
}
