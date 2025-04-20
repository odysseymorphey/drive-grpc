package utils

import "path/filepath"

func IsImage(filename string) bool {
	extensions := map[string]struct{}{
		".jpeg": {},
		".jpg":  {},
		".png":  {},
		".bmp":  {},
	}

	ext := filepath.Ext(filename)

	if _, ok := extensions[ext]; ok {
		return true
	}

	return false
}
