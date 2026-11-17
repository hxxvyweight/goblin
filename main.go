package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	incus "github.com/lxc/incus/client"
)

// Client holds: running status of program, working directory for
// incus spawns, number of active containers.
type Client struct {
	Status        bool
	Dir           string
	InstanceCount int
	Version       string
	Server        incus.InstanceServer
}

// Read log or text file and print individual seen file paths
// If specific files; e.g zip, tar, gz, iso, xz requiring
// extra steps are found, spawn new incus container, push the
// file(s) into the container and execute the standard reading
// process after unzipping/mounting/extracting/etc.
//
// Then output content of all files within the extracted folder
// into a file, prompt user then pull back to host if accepted.
func (c Client) readLog(input string) error {
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
			if c.hasGZSuffix(part) {
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
	c.printSlice(paths)
	return nil
}

// Iterate over slice and print
func (c Client) printSlice(validPaths []string) {
	for _, path := range validPaths {
		fmt.Printf("[+]: %s\n", path)
	}
}

// Check for .gz extension
func (c Client) hasGZSuffix(part string) bool {
	return strings.HasSuffix(part, ".gz")
}
func newClient() (*Client, time.Time, incus.InstanceServer, error) {

	t := time.Now()
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %v\n", err)
		return nil, t, nil, err
	}

	c := &Client{
		Status:        true,
		Dir:           dir,
		InstanceCount: 0,
		Version:       "goblin-0.1",
	}
	server := NewIncusInstance(c)
	c.Server = server

	fmt.Println(TermText)

	return c, t, server, nil
}

func endClient(t time.Time) {

	fmt.Println("\n" + TermText + "\n")
	fmt.Printf("|:: Elapsed time: %v\n", time.Since(t))
	fmt.Println("|:: Done" + "\n")
	ByeGoblin()
	fmt.Println("/|\\ Exiting...\n")
	os.Exit(1)

}
func main() {

	//New reader for Stdin
	//New goblin client
	reader := bufio.NewReader(os.Stdin)
	c, startTime, server, err := newClient()
	if err != nil {
		fmt.Printf("Error starting client: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("|> Started: %s\n|> Status: %t\n|> Dir: %s\n|> Containers: %d\n|> Version: [%s]\n", startTime, c.Status, c.Dir, c.InstanceCount, c.Version)
	fmt.Println(TermText)
	PrintGoblin()
	for c.Status {

		fmt.Print("\\|/ goblin |:: ")
		input, _ := reader.ReadString('\n')
		switch strings.TrimSpace(input) {
		case "1":
			fmt.Print("\\|/ goblin |:: Specify a file: ")
			input, _ = reader.ReadString('\n')
			c.readLog(strings.TrimSpace(input))
		case "create":
			CreateRequest(c, server)
		case "q":
			DeleteRequest(server)
			endClient(startTime)
		case "help":
			PrintHelp()
		default:
			continue
		}
	}

}
