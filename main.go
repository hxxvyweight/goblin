package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	incus "github.com/lxc/incus/client"
)

// * GoblinClient holds:
// - running status of program,
// - working directory for incus execution,
// - instance count,
// - version of program,
// - server instance ( see: github.com/lxc/incus/client#InstanceServer)

// This will be used and passed to the incus.go functions
// for API execution, the count is handled within each
// create and delete function: CreateRequest and DeleteRequest
// The server is received with newIncusInstance() and passed
// into the client.

// Status is here as another safety check if standard errors
// or unexpected occurrences happen.

type GoblinClient struct {
	Status        bool
	Dir           string
	ImageName     string
	InstanceCount int
	Version       string
	Server        incus.InstanceServer
}

// Print all current instances whether stopped or actively running
// Make a request through the incus api to receive a slice of []api.Instances
// Iterate and print
func (goblin GoblinClient) printInstances() {
	instances := GetInstancesRequest(goblin.Server)
	for _, in := range instances {
		fmt.Printf("Name: %s\nStatus: %v\nType: %v\n", in.Name, in.Status, in.Type)
	}
}

// Print paths of found files after viewing directory items
func (goblin GoblinClient) printPaths(validPaths []string) {
	for _, path := range validPaths {
		fmt.Printf("[+]: %s\n", path)
	}
}

// Check for .gz suffix
func (goblin GoblinClient) hasGZSuffix(part string) bool {
	return strings.HasSuffix(part, ".gz")
}

// Initialise goblin client, for communication with the incus API
// Start time of program is created
// Return to main
func newGoblinClient() (*GoblinClient, time.Time, incus.InstanceServer, error) {

	t := time.Now()
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %v\n", err)
		return nil, t, nil, err
	}

	g := &GoblinClient{
		Status:        true,
		Dir:           dir,
		ImageName:     "alpine",
		InstanceCount: 0,
		Version:       "goblin-0.1",
	}
	server, err := NewIncusInstance(g)
	if err != nil {
		fmt.Printf("error creating new incus server")
		log.Fatal(err)
	}
	g.Server = server

	fmt.Println(TermText)

	return g, t, server, nil
}

// End the client, print ASCII and terminal text
// Print elapsed time since start
func endClient(t time.Time) {

	fmt.Println("\n" + TermText + "\n")
	fmt.Printf("|:: Elapsed time: %v\n", time.Since(t))
	fmt.Println("|:: Done" + "\n")
	ByeGoblin()

	fmt.Println("/|\\ Exiting...\n")
	os.Exit(0)
}
func main() {

	//REPL for client
	//Initialise a new goblin client, receive start time of program
	//A server instance and a goblin client instance
	reader := bufio.NewReader(os.Stdin)
	goblin, startTime, server, err := newGoblinClient()
	if err != nil {
		fmt.Printf("Error starting client: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("|> Started: %s\n|> Status: %t\n|> Dir: %s\n|> Containers: %d\n|> Version: [%s]\n", startTime, goblin.Status, goblin.Dir, goblin.InstanceCount, goblin.Version)
	fmt.Println(TermText)
	PrintGoblin()
	for goblin.Status {
		fmt.Print("\\|/ goblin |:: ")
		input, _ := reader.ReadString('\n')

		switch strings.TrimSpace(input) {
		case "create":

			err = CreateRequest(goblin, server)
			if err != nil {
				fmt.Printf("error creating instance", err)
			}
		case "start":

			PrintInstancePrompt()
			input, _ = reader.ReadString('\n')
			err = StartRequest(goblin, goblin.Server, strings.TrimSpace(input))
			if err != nil {
				fmt.Printf("Failed to start instance:", err)

			}
		case "push":
			fileName := "test.txt"
			file := openFile(fileName)
			defer file.Close()

			PrintInstancePrompt()
			input, _ = reader.ReadString('\n')
			instanceNum := strings.TrimSpace(input)

			err = PushFileToWorker(goblin, instanceNum, file, fileName)
			if err != nil {
				fmt.Errorf("Error pushing file to worker: %v\n", err)
			}
		case "pull":
			PrintInstancePrompt()
			input, _ = reader.ReadString('\n')
			instanceNum := strings.TrimSpace(input)
			filePath := "/tmp/test.txt"

			err = PullFileFromWorker(goblin, instanceNum, filePath)
			if err != nil {
				fmt.Errorf("Error pulling file from worker: %v\n", err)
			}
		case "stop":
			PrintInstancePrompt()
			input, _ = reader.ReadString('\n')
			err = StopRequest(goblin, goblin.Server, strings.TrimSpace(input))
			if err != nil {
				fmt.Printf("Failed to stop instance:", err)

			}
		case "list":
			goblin.printInstances()
		case "delete":
			goblin.printInstances()
			PrintInstancePrompt()
			input, _ = reader.ReadString('\n')
			err = DeleteRequest(goblin, server, strings.TrimSpace(input))
			if err != nil {
				fmt.Printf("Error in deleting instance: %v\n", err)
				continue
			}
		case "q":
			endClient(startTime)
		case "help":
			PrintHelp()
		default:
			continue
		}
	}
}

func openFile(fileName string) *os.File {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	return f

}
