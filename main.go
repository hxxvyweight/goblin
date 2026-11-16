package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

const termText = "===================================================================="

// Client holds: running status of program, working directory for
// incus spawns, number of active containers
type Client struct {
	Status     bool
	Dir        string
	Containers int
}

// Read log or text file and print individual seen lines and total count
func readLog(reader *bufio.Reader) {
	seen := make(map[string]struct{})

	fmt.Print("\\|/ goblin |:: Specify a target file\n")
	fmt.Print("\\|/ goblin |:> ")
	input, err := reader.ReadString('\n')
	for {
		if input == "" || input == "\n" {
			continue
		}
		break
	}
	file, err := os.Open(strings.TrimSpace(input))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// New scanner for text file then
	// iterate and check if individual string in line contains specified extension
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.FieldsFunc(line, func(r rune) bool {
			return r == '/' || r == ' '

		})

		//Iterate over each string
		//Pass to helper func for .gz check
		for _, part := range parts {
			if hasGZSuffix(part) {
				seen[part] = struct{}{}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	//Make new slice to sort and return valid paths
	paths := make([]string, 0, len(seen))
	for path := range seen {
		paths = append(paths, path)
	}

	sort.Strings(paths)
	printSlice(paths)
}

func printSlice(validPaths []string) {
	for _, path := range validPaths {
		fmt.Printf("[+]: %s\n", path)
	}
}

func hasGZSuffix(part string) bool {
	return strings.HasSuffix(part, ".gz")
}

func main() {

	reader := bufio.NewReader(os.Stdin)

	c, startTime, err := newClient()
	if err != nil {
		fmt.Printf("Error starting client: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("|> Started: %s\n|> Status: %t\n|> Dir: %s\n|> Containers: %d\n", startTime, c.Status, c.Dir, c.Containers)
	fmt.Println(termText)

	printGoblin()

	fmt.Print("\\|/ goblin |:> ")
	for {
		input, _ := reader.ReadString('\n')
		switch strings.TrimSpace(input) {
		case "readlog":
			readLog(reader)
		case "q":
			fmt.Println("/|\\ Exiting...")
			os.Exit(1)
		default:
			continue
		}
	}

}
func execRead(c *Client, file *os.File, t time.Time) {
	SpawnContainer()
	endClient(t)
}

func newClient() (*Client, time.Time, error) {

	t := time.Now()
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %v\n", err)
		return nil, t, err
	}

	c := &Client{
		Status:     true,
		Dir:        dir,
		Containers: 0,
	}

	fmt.Println(termText)

	return c, t, nil
}

func endClient(t time.Time) {

	fmt.Println(termText + "\n")
	fmt.Printf("|> Elapsed time: %v\n", time.Since(t))
	fmt.Println("|> Done" + "\n")
	fmt.Println("/|\\ Exiting...")
	os.Exit(1)

}
func printGoblin() {

	fmt.Println(`
          /\_/\                    _______
         ( o.o )       .-.        /       \
         / > ^ \      /   \      |  $ $ $  |
        / / \ \ \    | L O O T!  |_________|
       (_/   \_)      \___/      [___________]
      /           \              | $$$$$$$$$ |
     /    GOBLIN   \             | $$$$$$$$$ |
    /_______________\            |___________|`)

	fmt.Println(">>> GOBLIN CLIENT INITIALIZED!")
	fmt.Println(">>> Ready to venture into the archives...\n")
}
