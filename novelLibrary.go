package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func GetSubfolders(root string) ([]string, error) {
	var folders []string
	files, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			folders = append(folders, file.Name())
		}
	}
	return folders, nil
}

func GetFiles(dir string) ([]string, error) {
	log.Printf("Getting files from %s", dir)
	files, err := os.ReadDir(dir)
	if err != nil {
		return make([]string, 0), err
	}

	var j = 0
	chapters := make([]string, len(files))
	for _, file := range files {
		match, _ := filepath.Match("*.txt", file.Name())
		if match {
			chapters[j] = file.Name()
			j++
		}
	}
	sort.Slice(chapters[0:j], func(x, y int) bool {
		numX, _ := extractNumber(chapters[x])
		numY, _ := extractNumber(chapters[y])
		return numX < numY
	})

	log.Printf("%s, %d", dir, len(chapters))

	return chapters[0:j], nil
}

func GetChapter(novel string, chapter string) (string, error) {
	filePath := CONTENT + "/" + novel + "/" + chapter
	file, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(file), nil
}

func SaveSelectedChapter(novel string, chapter string) error {
	filePath := CONTENT + "/" + novel + "/" + chapter
	saveGamePath := CONTENT + "/" + novel + "/" + "chapter.save"
	_, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = os.WriteFile(saveGamePath, []byte(chapter), 0644)
	if err != nil {
		return err
	}

	return nil
}

func GetSelectedChapter(novel string) (string, error) {
	saveGamePath := CONTENT + "/" + novel + "/" + "chapter.save"
	file, err := os.ReadFile(saveGamePath)
	if err != nil {
		return "EP.0.txt", err
	}
	chapter := string(file)
	return chapter, nil
}

func extractNumber(fileName string) (int, error) {
	parts := strings.Split(fileName, ".")
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid file name")
	}
	return strconv.Atoi(parts[1])
}

func GetPrependText(novel string) (string, error) {
	filePath := CONTENT + "/" + novel + "/prepend.prefix"
	file, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(file), nil
}

func SavePrependText(novel string, text string) error {
	filePath := CONTENT + "/" + novel + "/prepend.prefix"
	return os.WriteFile(filePath, []byte(text), 0644)
}
