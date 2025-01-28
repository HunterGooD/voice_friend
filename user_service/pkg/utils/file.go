package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func FileExist(basePath, filePath string) (bool, error) {
	fullPath, err := GetFullPath(basePath, filePath)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(fullPath)
	return err == nil, nil
}

func GetFullPath(basePath, filePath string) (string, error) {
	const NameFunc = "utils.GetFullPath"
	var err error
	if basePath == "" {
		basePath, err = os.Getwd()
		if err != nil {
			return "", fmt.Errorf("%s: failed to get current dir: %w", NameFunc, err)
		}
	}
	fullPath, err := filepath.Abs(filepath.Join(basePath, filePath))
	if err != nil {
		return "", fmt.Errorf("%s: failed to create absolute path: %w", NameFunc, err)
	}
	return fullPath, err
}
