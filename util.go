package main

import (
	"errors"
	"fmt"
	"os"
)

func getFileInfo(filename string) (os.FileInfo, error) {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return nil, err
	}
	return info, nil
}

func IsFile(filename string) bool {
	info, err := getFileInfo(filename)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return !info.IsDir()
}

func IsDir(filename string) bool {
	info, err := getFileInfo(filename)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return info.IsDir()
}

func ImagesDir(image string) string {
	return fmt.Sprintf("%s/%s", *RepoRoot, REPO_IMAGES_DIR)
}

func ImageLabelDir(image string) string {
	return fmt.Sprintf("%s/%s", ImagesDir(image), image)
}

func ImageTagDir(image, tag string) string {
	return fmt.Sprintf("%s/%s/%s", ImageLabelDir(image), IMAGE_TAGS_DIR_POSTFIX, tag)
}