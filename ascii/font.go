package ascii

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func getBaseDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		return filepath.Dir(filepath.Dir(filename))
	}
	return "."
}

func LoadBanner(bannerType string) ([]string, error) {
	bannerName := bannerType
	if !strings.HasSuffix(bannerName, ".txt") {
		bannerName += ".txt"
	}

	templatePath := filepath.Join(getBaseDir(), "template", bannerName)

	data, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("error reading banner file %s: %v", bannerName, err)
	}

	content := strings.ReplaceAll(string(data), "\r\n", "\n")

	return strings.Split(content, "\n"), nil
}
