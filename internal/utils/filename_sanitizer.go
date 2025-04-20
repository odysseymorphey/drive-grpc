package utils

import (
	"path/filepath"
	"strings"
)

func SanitizeFilename(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}

	name = filepath.Base(name)
	name = strings.ReplaceAll(name, string(filepath.Separator), "")

	return strings.ReplaceAll(name, "\x00", "")
}
