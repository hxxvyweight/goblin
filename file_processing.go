package main

import (
	"fmt"
	"os"
)

const INPUT_DIR = "./queue"

func GetFileNames() []os.DirEntry {
	createQueueDir()

	entries, err := os.ReadDir(INPUT_DIR)
	if err != nil {
		fmt.Errorf("Error creating input directory", err)
	}

	if len(entries) < 1 {
		fmt.Errorf("No directory contents:", err)
	}
	return entries

}

func PrintFilePaths() {
	entries := GetFileNames()
	for i, e := range entries {
		fmt.Printf("File no: %d: %s\n", i, e.Name())
	}
}

func createQueueDir() {
	if err := os.MkdirAll(INPUT_DIR, 0755); err != nil {
		fmt.Errorf("Error creating directory ./goblin/queue for files intake", err)
	}
	fmt.Println("Directory created or already exists:", INPUT_DIR)
}
