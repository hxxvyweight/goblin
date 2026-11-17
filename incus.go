package main

import (
	"fmt"
	"log"
	"os"

	incus "github.com/lxc/incus/client"
	"github.com/lxc/incus/shared/api"
)

var instanceName string

func NewIncusInstance(g *GoblinClient) (u incus.InstanceServer) {
	u, err := incus.ConnectIncusUnix("", nil)
	if err != nil {
		log.Fatalf("Failed to connect to Incus: %v", err)
	}

	return u
}

func GetInstancesRequest(u incus.InstanceServer) []api.Instance {
	instances, err := u.GetInstances(api.InstanceTypeContainer)
	if err != nil {
		fmt.Printf("Failed to get instances: %v\n", err)
	}
	return instances
}

//API Execution
//CreateRequest(u, instanceName)
//StartRequest(u, instanceName)
//ExecRequest(u, instanceName)
//StopRequest(u, instanceName)
//DeleteRequest(u, instanceName)

func CreateRequest(goblin *GoblinClient, u incus.InstanceServer) {
	goblin.InstanceCount++
	instanceName = fmt.Sprintf("worker-container-%d", goblin.InstanceCount)
	fmt.Printf("Instances: %d\n", goblin.InstanceCount)
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
}

func ExecRequest(u incus.InstanceServer, instanceName string) {
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

func StopRequest(u incus.InstanceServer, instanceName string) {

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

func StartRequest(u incus.InstanceServer, instanceName string) {

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

func DeleteRequest(g *GoblinClient, u incus.InstanceServer, instanceName string) {

	fmt.Printf("Deleting instance: %s\n", instanceName)
	deleteOp, err := u.DeleteInstance(instanceName)
	if err != nil {
		log.Fatal(err)
	}
	if err = deleteOp.Wait(); err != nil {
		log.Fatal(err)
	}
}
