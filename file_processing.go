package main

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strings"
)

// Read log or text file and print individual seen file paths
// If specific files; e.g zip, tar, gz, iso, xz requiring
// extra steps are found, spawn new incus container, push the
// file(s) into the container and execute the standard reading
// process after unzipping/mounting/extracting/etc.
//
// Then output content of all files within the extracted folder
// into a file, prompt user then pull back to host if accepted.
func ReadLog(input string) error {
	seen := make(map[string]struct{})

	file, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// New scanner for text file then
	// iterate and match individual strings to extensions
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.FieldsFunc(line, func(r rune) bool {
			return r == '/' || r == ' '

		})

		//Iterate over each string
		//Pass to helper func for .gz check
		for _, part := range parts {
			if goblin.hasGZSuffix(part) {
				seen[part] = struct{}{}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	//Make new slice to sort and print valid paths
	paths := make([]string, 0, len(seen))
	for path := range seen {
		paths = append(paths, path)
	}

	sort.Strings(paths)
	return paths
}
