package ascii

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// GetTerminalWidth returns the current terminal width
func GetTerminalWidth() int {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return DefaultWidth
	}

	parts := strings.Fields(strings.TrimSpace(string(out)))
	if len(parts) < 2 {
		return DefaultWidth
	}

	w, err := strconv.Atoi(parts[1])
	if err != nil || w <= 0 {
		return DefaultWidth
	}

	return w
}
