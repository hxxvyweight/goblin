package main

import (
	"fmt"
	"log"
	"os"

	incus "github.com/lxc/incus/client"
	"github.com/lxc/incus/shared/api"
)

func SpawnContainer() {
	u, err := incus.ConnectIncusUnix("", nil)
	if err != nil {
		log.Fatalf("Failed to connect to Incus: %v", err)
	}

	instanceName := "worker-container-1"

	req := api.InstancesPost{
		Name: instanceName,
		Type: api.InstanceTypeContainer,
		Source: api.InstanceSource{
			Type:  "image",
			Alias: "alpine-worker",
		},
		Config: map[string]string{
			"security.privileged": "true",
			"security.nesting":    "true",
			"raw.lxc":             "lxc.cgroup.dir = incus\nlxc.mount.auto = cgroup:mixed",
		},
	}

	op, err := u.CreateInstance(req)
	if err != nil {
		log.Printf("Instance might already exist, proceeding to start... (%v)", err)
	} else {
		if err = op.Wait(); err != nil {
			log.Fatalf("Instance creation operation failed: %v", err)
		}
		log.Printf("Successfully created container instance: %s", instanceName)
	}

	log.Printf("Container %s is now running", instanceName)

	//API Execution
	startRequest(u, instanceName)
	execRequest(u, instanceName)
	stopRequest(u, instanceName)
	deleteRequest(u, instanceName)

}

func execRequest(u incus.InstanceServer, instanceName string) {
	execReq := api.InstanceExecPost{
		Command:     []string{"ls", "-la"},
		WaitForWS:   true,
		Interactive: false,
	}

	execArgs := incus.InstanceExecArgs{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	execOp, err := u.ExecInstance(instanceName, execReq, &execArgs)
	if err != nil {
		log.Fatal(err)
	}

	if err = execOp.Wait(); err != nil {
		log.Fatal(err)
	}
}

func stopRequest(u incus.InstanceServer, instanceName string) {

	stopReq := api.InstanceStatePut{
		Action:  "stop",
		Timeout: -1,
	}

	fmt.Printf("Stopping instance: %s\n", instanceName)

	stateOp, err := u.UpdateInstanceState(instanceName, stopReq, "")
	if err != nil {
		log.Fatal(err)
	}

	if err = stateOp.Wait(); err != nil {
		log.Fatal(err)
	}

}

func startRequest(u incus.InstanceServer, instanceName string) {

	startReq := api.InstanceStatePut{
		Action:  "start",
		Timeout: -1,
	}
	stateOp, err := u.UpdateInstanceState(instanceName, startReq, "")
	if err != nil {
		log.Fatalf("Failed to trigger start operation: %v", err)
	}

	if err = stateOp.Wait(); err != nil {
		log.Fatalf("Container failed to start: %v", err)
	}
}

func deleteRequest(u incus.InstanceServer, instanceName string) {

	fmt.Printf("Deleting instance: %s\n", instanceName)

	deleteOp, err := u.DeleteInstance("worker-container-1")
	if err != nil {
		log.Fatal(err)
	}
	if err = deleteOp.Wait(); err != nil {
		log.Fatal(err)
	}
}
