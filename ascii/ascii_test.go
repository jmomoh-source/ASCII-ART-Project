package ascii

import (
	"testing"
)

func TestAsciiArt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single character",
			input:    "A",
			expected: "           \n    /\\     \n   /  \\    \n  / /\\ \\   \n / ____ \\  \n/_/    \\_\\ \n           \n           \n",
		},
		{
			name:     "hello lowercase",
			input:    "hello",
			expected: " _              _   _          \n| |            | | | |         \n| |__     ___  | | | |   ___   \n|  _ \\   / _ \\ | | | |  / _ \\  \n| | | | |  __/ | | | | | (_) | \n|_| |_|  \\___| |_| |_|  \\___/  \n                               \n                               \n",
		},
		{
			name:     "HeLlO mixed case",
			input:    "HeLlO",
			expected: " _    _          _        _    ____   \n| |  | |        | |      | |  / __ \\  \n| |__| |   ___  | |      | | | |  | | \n|  __  |  / _ \\ | |      | | | |  | | \n| |  | | |  __/ | |____  | | | |__| | \n|_|  |_|  \\___| |______| |_|  \\____/  \n                                      \n                                      \n",
		},
		{
			name:     "newline split",
			input:    "Hi\\nThere",
			expected: " _    _   _  \n| |  | | (_) \n| |__| |  _  \n|  __  | | | \n| |  | | | | \n|_|  |_| |_| \n             \n             \n _______   _                           \n|__   __| | |                          \n   | |    | |__     ___   _ __    ___  \n   | |    |  _ \\   / _ \\ | '__|  / _ \\ \n   | |    | | | | |  __/ | |    |  __/ \n   |_|    |_| |_|  \\___| |_|     \\___| \n                                       \n                                       \n",
		},
		{
			name:     "double newline",
			input:    "Hi\\n\\nThere",
			expected: " _    _   _  \n| |  | | (_) \n| |__| |  _  \n|  __  | | | \n| |  | | | | \n|_|  |_| |_| \n             \n             \n\n _______   _                           \n|__   __| | |                          \n   | |    | |__     ___   _ __    ___  \n   | |    |  _ \\   / _ \\ | '__|  / _ \\ \n   | |    | | | | |  __/ | |    |  __/ \n   |_|    |_| |_|  \\___| |_|     \\___| \n                                       \n                                       \n",
		},
		{
			name:     "space character",
			input:    " ",
			expected: "      \n      \n      \n      \n      \n      \n      \n      \n",
		},
		{
			name:     "exclamation mark",
			input:    "!",
			expected: " _  \n| | \n| | \n| | \n|_| \n(_) \n    \n    \n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AsciiArt(tt.input)
			if result != tt.expected {
				t.Errorf("AsciiArt(%q):\nGot:\n%s\nExpected:\n%s", tt.input, result, tt.expected)
			}
		})
	}
}
