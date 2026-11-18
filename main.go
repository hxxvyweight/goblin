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
func (goblin *GoblinClient) printInstances() {
	instances := GetInstancesRequest(goblin)
	for _, in := range instances {
		fmt.Printf("Name: %s\nStatus: %v\nType: %v\n", in.Name, in.Status, in.Type)
	}
}

// Initialise goblin client, for communication with the incus API
// Start time of program is created
// Return to main
func newGoblinClient() (*GoblinClient, time.Time, error) {

	t := time.Now()
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %v\n", err)
		return nil, t, err
	}

	g := &GoblinClient{
		Status:        true,
		Dir:           dir,
		ImageName:     "alpine",
		InstanceCount: 0,
		Version:       "goblin-0.1",
	}
	server, err := NewIncusInstance()
	if err != nil {
		fmt.Printf("error creating new incus server: %v\n", err)
		return nil, t, err
	}
	g.Server = server

	fmt.Println(TermText)

	return g, t, nil
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
	goblin, startTime, err := newGoblinClient()
	if err != nil {
		fmt.Printf("Error starting client: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("|> Started: %s\n|> Status: %t\n|> Dir: %s\n|> Containers: %d\n|> Version: [%s]\n", startTime, goblin.Status, goblin.Dir, goblin.InstanceCount, goblin.Version)
	fmt.Println(TermText)
	PrintGoblin()

	for goblin.Status {

		fmt.Print("\\|/ goblin |:: ")
		input := getInput(reader)

		switch input {
		case "create":
			create(goblin, reader)
		case "delete":
			delete(goblin, reader)
		case "delete all":
			deleteAll(goblin)
		case "start":
			start(goblin, reader)
		case "push":
			push(goblin, reader)
		case "pull":
			err := ExecRequest(goblin, "1")
			if err != nil {
				fmt.Println(err)
			}
			pull(goblin, reader)
		case "stop":
			stop(goblin, reader)
		case "queue":
			PrintFilePaths()
		case "list":
			goblin.printInstances()
		case "exit":
			endClient(startTime)
		case "help":
			PrintHelp()
		default:
			fmt.Println("Unknown command. Type 'help'.")
		}
	}
}

func getInput(reader *bufio.Reader) string {
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading user input: %v\n", err)
			return ""
		}
		if s != "" {
			return strings.TrimSpace(s)
		}
	}

}

// Create, Delete, Stop, Start, Push, Pull
func pull(goblin *GoblinClient, reader *bufio.Reader) {
	PrintInstancePrompt()
	instanceNum := getInput(reader)
	//entries := GetFileNames()
	reportPath := fmt.Sprintf("/tmp/report.json")
	err := PullFileFromWorker(goblin, instanceNum, reportPath)
	if err != nil {
		fmt.Printf("Error pulling file from worker: %v | File name: report.json\n", err)
	}
}
func push(goblin *GoblinClient, reader *bufio.Reader) {
	PrintInstancePrompt()
	instanceNum := getInput(reader)
	fmt.Println("Getting entries")
	entries := GetFileNames()
	for i, entry := range entries {

		fmt.Printf("Opening file [%d]: %s\n", i, entry.Name())
		filePath := fmt.Sprintf("queue/%s", entry.Name())

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			return
		}
		err = PushFileToWorker(goblin, instanceNum, file, entry.Name())
		if err != nil {
			fmt.Printf("Error pushing file to worker: %v | File no: %d\nFile name: %s\n", err, i, entry.Name())
		}
		file.Close()
	}
	filePath := fmt.Sprintf("workers/analyzer-zip/target/x86_64-unknown-linux-musl/release/analyzer-zip")
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening Rust executable: %v\n", err)
	}
	err = PushFileToWorker(goblin, instanceNum, file, "analyzer-zip")
	file.Close()
	if err != nil {
		fmt.Printf("Error pushing file to worker: %v", err)
	}

}
func create(goblin *GoblinClient, reader *bufio.Reader) {
	err := CreateRequest(goblin)
	if err != nil {
		fmt.Printf("failed to create instance: %v\n", err)
	}
	goblin.printInstances()
}
func start(goblin *GoblinClient, reader *bufio.Reader) {
	PrintInstancePrompt()
	input := getInput(reader)
	err := StartRequest(goblin, input)
	if err != nil {
		fmt.Printf("failed to start instance:", err)

	}
}
func stop(goblin *GoblinClient, reader *bufio.Reader) {

	PrintInstancePrompt()
	input := getInput(reader)
	err := StopRequest(goblin, input)
	if err != nil {
		fmt.Printf("failed to stop instance:", err)

	}
}
func delete(goblin *GoblinClient, reader *bufio.Reader) {
	goblin.printInstances()
	PrintInstancePrompt()
	input := getInput(reader)
	err := DeleteRequest(goblin, input)
	if err != nil {
		fmt.Printf("failed to delete instance: %v\n", err)
	}
}
func deleteAll(goblin *GoblinClient) {
	instances := GetInstancesRequest(goblin)
	for _, entry := range instances {
		if entry.Status != "Stopped" {
			err := StopAllRequest(goblin, entry.Name)
			if err != nil {
				fmt.Printf("Error stopping instance: %s\n", entry.Name)
				continue
			}
		}
		err := DeleteAllRequest(goblin, entry.Name)
		if err != nil {
			fmt.Printf("Error deleting entry: %v\n", entry)
		}
	}

}
