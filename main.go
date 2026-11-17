package main

import (
	"bufio"
	"fmt"
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
	InstanceCount int
	Version       string
	Server        incus.InstanceServer
}

// Print all current instances whether stopped or actively running
// Make a request through the incus api to receive a slice of []api.Instances
// Iterate and print
func (goblin GoblinClient) printInstances(server incus.InstanceServer) {
	instances := GetInstancesRequest(server)
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

// Return a instance of the goblin client with start time of the program
// Return a instance of the incus server
// Return error if working dir not found
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
		InstanceCount: 0,
		Version:       "goblin-0.1",
	}
	server := NewIncusInstance(g)
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
	//Initiate a new goblin client, receive start time of program
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
		if goblin.Status {

			fmt.Print("\\|/ goblin |:: ")
			input, _ := reader.ReadString('\n')

			switch strings.TrimSpace(input) {
			case "create":
				CreateRequest(goblin, server)
			case "list":
				goblin.printInstances(server)
			case "delete":
				fmt.Print("\\|/ goblin |:: Specify a instance: ")
				input, _ = reader.ReadString('\n')
				DeleteRequest(goblin, server, strings.TrimSpace(input))
				goblin.printInstances(server)
			case "q":
				endClient(startTime)
			case "help":
				PrintHelp()
			default:
				continue
			}
		} else {
			fmt.Printf("Error - client not running: %t\n", goblin.Status)
			os.Exit(1)
		}
	}

}
