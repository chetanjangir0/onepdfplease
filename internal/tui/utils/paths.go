package utils

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func GetNextAvailablePath(basePath string) string {
	if _, err := os.Stat(basePath); err != nil {
		return basePath
	}

	ext := filepath.Ext(basePath)
	nameWithoutExt := strings.TrimSuffix(basePath, ext)

	counter := 1
	for {
		newPath := fmt.Sprintf("%s_%d%s", nameWithoutExt, counter, ext)

		if _, err := os.Stat(newPath); errors.Is(err, fs.ErrNotExist) {
			return newPath
		}
		counter++
	}
}

func UpdateFileExtension(path, newExt string) string {
	ext := filepath.Ext(path)
	switch ext {
	case newExt:
		return path
	case "":
		return path + newExt
	}

	nameWithoutExt := strings.TrimSuffix(path, ext)
	return nameWithoutExt + newExt
}
