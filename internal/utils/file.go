package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func LocateNearestFile(targetFile string) (src string, err error) {
	currentPath, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dirs := strings.Split(currentPath, "/")

	for i := len(dirs); i > 0; i-- {
		searchableDir := "/"
		for j := 0; j < i; j++ {
			searchableDir = filepath.Join(searchableDir, dirs[j])
		}

		fullSearchedFile := filepath.Join(searchableDir, targetFile)

		info, err := os.Stat(fullSearchedFile)
		if os.IsNotExist(err) {
			continue
		}

		if !info.IsDir() {
			return fullSearchedFile, nil
		}
	}

	return "", fmt.Errorf("directory %s not found", targetFile)
}
