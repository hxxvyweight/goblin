package main

import (
	"fmt"
	"log"
	"os"

	incus "github.com/lxc/incus/client"
	"github.com/lxc/incus/shared/api"
)

var instanceName string

func NewIncusInstance(g *GoblinClient) (u incus.InstanceServer, err error) {
	u, err = incus.ConnectIncusUnix("", nil)
	if err != nil {
		log.Fatalf("Failed to connect to Incus: %v", err)
		return nil, err
	}

	return u, nil
}

func GetInstancesRequest(u incus.InstanceServer) []api.Instance {
	instances, err := u.GetInstances(api.InstanceTypeContainer)
	if err != nil {
		fmt.Printf("Failed to get instances: %v\n", err)
	}
	return instances
}

//API Execution
//CreateRequest(goblin, server)
//StartRequest(goblin, server)
//ExecRequest(goblin, server)
//StopRequest(goblin, server)
//DeleteRequest(goblin, server)

func CreateRequest(goblin *GoblinClient, u incus.InstanceServer) error {
	goblin.InstanceCount++
	instanceName = fmt.Sprintf("worker-%s-%d", goblin.ImageName, goblin.InstanceCount)
	fmt.Printf("Instances: %d\n", goblin.InstanceCount)
	req := api.InstancesPost{
		Name: instanceName,
		Type: api.InstanceTypeContainer,
		Source: api.InstanceSource{
			Type:  "image",
			Alias: "alpine-worker",
		},
		Config: map[string]string{
			"security.privileged": "false",
			"security.nesting":    "true",
		},
	}

	op, err := u.CreateInstance(req)
	if err != nil {
		fmt.Printf("Error creating instance operation: %s\n%v", instanceName, err)
		return err
	} else {
		if err = op.Wait(); err != nil {
			fmt.Printf("Instance creation failed: %s\n%v\n", instanceName, err)
			return err
		}
		log.Printf("Successfully created container instance: %s", instanceName)
	}
	return nil
}

func ExecRequest(u incus.InstanceServer, instanceName string) error {
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
		fmt.Printf("Failed to perform exec operation on instance: %s\n%v\n", instanceName, err)
		return err

	}

	if err = execOp.Wait(); err != nil {
		fmt.Printf("Error on exec Wait() operation: %s\n%v\n", instanceName, err)
		return err
	}
	return nil
}

func StopRequest(u incus.InstanceServer, instanceName string) error {

	stopReq := api.InstanceStatePut{
		Action:  "stop",
		Timeout: -1,
	}

	fmt.Printf("Stopping instance: %s\n", instanceName)

	stateOp, err := u.UpdateInstanceState(instanceName, stopReq, "")
	if err != nil {
		fmt.Printf("Failed to perform stop operation on instance: %s\n%v\n", instanceName, err)
		return err
	}

	if err = stateOp.Wait(); err != nil {
		fmt.Printf("Error on stop Wait() operation: %s\n%v\n", instanceName, err)
		return err
	}
	return nil

}

func StartRequest(u incus.InstanceServer, instanceName string) error {

	startReq := api.InstanceStatePut{
		Action:  "start",
		Timeout: -1,
	}
	stateOp, err := u.UpdateInstanceState(instanceName, startReq, "")
	if err != nil {
		fmt.Printf("Failed to perform start operation on instance: %s\n%v\n", instanceName, err)
		return err
	}

	if err = stateOp.Wait(); err != nil {
		fmt.Printf("Error on start Wait() operation: %v\n", err)
		return err
	}
	return nil
}

func DeleteRequest(g *GoblinClient, u incus.InstanceServer, instanceNum string) error {

	fmt.Printf("Deleting instance: worker-%s-%s\n", g.ImageName, instanceNum)
	deleteOp, err := u.DeleteInstance(instanceName)
	if err != nil {
		fmt.Printf("Failed to perform delete operation on instance: worker-%s-%s\n%v\n", instanceName, instanceNum, err)
		return err
	}
	if err = deleteOp.Wait(); err != nil {
		fmt.Printf("Error on delete Wait() operation: %s\n%v\n", instanceName, err)
		return err
	}
	return nil
}
