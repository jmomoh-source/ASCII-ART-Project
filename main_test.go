package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"ascii-art/ascii"
)

// runMain captures the output of the main() function by overriding os.Args and os.Stdout
func runMain(args []string) string {
	oldArgs := os.Args
	oldStdout := os.Stdout
	defer func() {
		os.Args = oldArgs
		os.Stdout = oldStdout
	}()

	// Set os.Args. The first argument is usually the program name, so we add a dummy "cmd"
	os.Args = append([]string{"cmd"}, args...)

	// Create a pipe to capture the output of fmt.Print / fmt.Println
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call main directly
	main()

	// Close the writer and restore os.Stdout so we can read from the pipe
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// --- COLOR TESTS ---

func TestInvalidFlagFormat(t *testing.T) {
	output := runMain([]string{"--color", "red", "banana"})
	expectedUsage := "Usage: go run . [OPTION] [STRING] [BANNER]"

	if !strings.Contains(output, expectedUsage) {
		t.Errorf("Expected usage message, got:\n%s", output)
	}
}

func TestWholeStringColored(t *testing.T) {
	output := runMain([]string{"--color=red", "hello world"})
	redANSI := "\033[31m"

	if !strings.Contains(output, redANSI) {
		t.Errorf("Expected red ANSI code %q in output, but was missing", redANSI)
	}
	if strings.Contains(output, "Usage:") {
		t.Errorf("Unexpected usage message in output:\n%s", output)
	}
}

func TestSpecificSubstringColored(t *testing.T) {
	output := runMain([]string{"--color=orange", "GuYs", "HeY GuYs"})
	orangeANSI := "\033[38;2;255;165;0m" 

	if !strings.Contains(output, orangeANSI) {
		t.Errorf("Expected orange ANSI code %q in output, but was missing", orangeANSI)
	}
}

func TestSingleLetterColored(t *testing.T) {
	output := runMain([]string{"--color=blue", "B", "RGB()"})
	blueANSI := "\033[34m"

	if !strings.Contains(output, blueANSI) {
		t.Errorf("Expected blue ANSI code %q in output, but was missing", blueANSI)
	}
}

func TestMultipleColorFlags(t *testing.T) {
	output := runMain([]string{"--color=red", "A", "--color=blue", "B", "A B C"})
	redANSI := "\033[31m"
	blueANSI := "\033[34m"

	if !strings.Contains(output, redANSI) {
		t.Errorf("Expected red ANSI code %q in output, but was missing", redANSI)
	}
	if !strings.Contains(output, blueANSI) {
		t.Errorf("Expected blue ANSI code %q in output, but was missing", blueANSI)
	}
}

func TestHSLColor(t *testing.T) {
	output := runMain([]string{"--color=hsl(0, 100%, 50%)", "red"})
	// HSL(0, 100%, 50%) maps perfectly to red RGB: 255, 0, 0
	expectedANSI := "\033[38;2;255;0;0m"
	
	if !strings.Contains(output, expectedANSI) {
		t.Errorf("Expected HSL calculated ANSI code %q for red, but got missing", expectedANSI)
	}
}

// --- FS / FONT TESTS ---

func TestBannerFilesExist(t *testing.T) {
	banners := []string{"standard", "shadow", "thinkertoy"}

	for _, banner := range banners {
		_, err := ascii.LoadBanner(banner)
		if err != nil {
			t.Errorf("Banner %s does not exist or failed to load: %v", banner, err)
		}
	}
}

// --- JUSTIFY / ALIGN TESTS ---

// TestLeftAlignment tests that the left alignment correctly returns the original string without padding
func TestLeftAlignment(t *testing.T) {
	line := "hello"
	width := 10
	expected := "hello" // Left alignment shouldn't add any spaces

	result := ascii.AlignLine(line, ascii.AlignLeft, width)
	if result != expected {
		t.Errorf("Expected: %q, but got %q", expected, result)
	}
}

// TestRightAlignment tests that the right alignment correctly calculates empty spaces to push the word to the right edge
func TestRightAlignment(t *testing.T) {
	line := "hello"
	width := 10
	expected := "     hello" // Terminal Width (10) - "hello" Length (5) = 5 spaces added before the word

	result := ascii.AlignLine(line, ascii.AlignRight, width)
	if result != expected {
		t.Errorf("Expected: %q, but got %q", expected, result)
	}
}

// TestCenterAlignment tests that the center alignment mathematical formula works
func TestCenterAlignment(t *testing.T) {
	line := "abc"
	width := 9
	expected := "   abc" // Terminal Width (9) - "abc" Length (3) = 6 spaces leftover. Division by 2 = 3 front spaces

	result := ascii.AlignLine(line, ascii.AlignCenter, width)
	if result != expected {
		t.Errorf("Expected: %q, but got %q", expected, result)
	}
}

// TestAlignJustify ensures spacing is cleanly spread when using justify align mode
func TestAlignJustifyIntegration(t *testing.T) {
	output := runMain([]string{"--align=justify", "test output"})
	// Without parsing exact spacing geometry mapping (as this tests terminal width logic),
	// We at least ensure the wrapper evaluates justify correctly outputting results without crashing.
	if output == "" || strings.Contains(output, "Usage:") {
		t.Errorf("Justify alignment failed, unexpectedly empty or errored.")
	}
}

// --- OUTPUT TESTS ---

func TestFileOutput(t *testing.T) {
	tempFileName := "temp_test_output.txt"
	defer os.Remove(tempFileName)

	output := runMain([]string{"--output=" + tempFileName, "hello", "standard"})
	
	// CLI shouldn't output to terminal when using --output
	if output != "" {
		t.Errorf("Expected empty terminal output when using --output, got: %v", output)
	}

	content, err := os.ReadFile(tempFileName)
	if err != nil {
		t.Fatalf("Failed to read expected output file: %v", err)
	}

	if len(content) == 0 {
		t.Errorf("Output file is completely empty.")
	}
}
