package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"ascii-art/ascii"
)

func runMain(args []string) string {
	oldArgs := os.Args
	oldStdout := os.Stdout

	defer func() {
		os.Args = oldArgs
		os.Stdout = oldStdout
	}()

	os.Args = append([]string{"cmd"}, args...)

	r, w, _ := os.Pipe()

	os.Stdout = w

	main()

	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}

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
	expectedANSI := "\033[38;2;255;0;0m"

	if !strings.Contains(output, expectedANSI) {
		t.Errorf("Expected HSL calculated ANSI code %q for red, but got missing", expectedANSI)
	}
}

func TestBannerFilesExist(t *testing.T) {
	banners := []string{"standard", "shadow", "thinkertoy"}

	for _, banner := range banners {
		_, err := ascii.LoadBanner(banner)
		if err != nil {
			t.Errorf("Banner %s does not exist or failed to load: %v", banner, err)
		}
	}
}

func TestLeftAlignment(t *testing.T) {
	line := "hello"
	width := 10
	expected := "hello"

	result := ascii.AlignLine(line, ascii.AlignLeft, width)
	if result != expected {
		t.Errorf("Expected: %q, but got %q", expected, result)
	}
}

func TestRightAlignment(t *testing.T) {
	line := "hello"
	width := 10
	expected := "     hello"

	result := ascii.AlignLine(line, ascii.AlignRight, width)
	if result != expected {
		t.Errorf("Expected: %q, but got %q", expected, result)
	}
}

func TestCenterAlignment(t *testing.T) {
	line := "abc"
	width := 9
	expected := "   abc"

	result := ascii.AlignLine(line, ascii.AlignCenter, width)
	if result != expected {
		t.Errorf("Expected: %q, but got %q", expected, result)
	}
}

func TestAlignJustifyIntegration(t *testing.T) {
	output := runMain([]string{"--align=justify", "test output"})
	if output == "" || strings.Contains(output, "Usage:") {
		t.Errorf("Justify alignment failed, unexpectedly empty or errored.")
	}
}

func TestFileOutput(t *testing.T) {
	tempFileName := "temp_test_output.txt"
	defer os.Remove(tempFileName)

	output := runMain([]string{"--output=" + tempFileName, "hello", "standard"})

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
