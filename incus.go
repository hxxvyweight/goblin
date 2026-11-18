package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	incus "github.com/lxc/incus/client"
	"github.com/lxc/incus/shared/api"
)

var instanceName string

func NewIncusInstance() (u incus.InstanceServer, err error) {
	u, err = incus.ConnectIncusUnix("", nil)
	if err != nil {
		log.Fatalf("Failed to connect to Incus: %v", err)
		return nil, err
	}

	return u, nil
}

func GetInstancesRequest(goblin *GoblinClient) []api.Instance {
	instances, err := goblin.Server.GetInstances(api.InstanceTypeContainer)
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

func CreateRequest(goblin *GoblinClient) error {
	goblin.InstanceCount++
	instanceNum := fmt.Sprintf("%d", goblin.InstanceCount)
	instanceName := fmt.Sprintf("worker-%s-%s", goblin.ImageName, instanceNum)
	req := api.InstancesPost{
		Name: instanceName,
		Type: api.InstanceTypeContainer,
		Source: api.InstanceSource{
			Type:  "image",
			Alias: "alpine-worker",
		},
		Config: map[string]string{
			"security.privileged": "false",
			"security.nesting":    "false",
			"raw.lxc":             "lxc.mount.auto = proc:mixed sys:ro cgroup:ro",
		},
	}

	op, err := goblin.Server.CreateInstance(req)
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

func ExecRequest(goblin *GoblinClient, instanceNum string) error {
	instanceName := fmt.Sprintf("worker-%s-%s", goblin.ImageName, instanceNum)
	fmt.Printf("Executing on instance: worker-%s-%s\n", goblin.ImageName, instanceNum)

	execReq := api.InstanceExecPost{
		Command:     []string{"ls", "-la"},
		WaitForWS:   true,
		Interactive: false,
	}

	execArgs := incus.InstanceExecArgs{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	execOp, err := goblin.Server.ExecInstance(instanceName, execReq, &execArgs)
	if err != nil {
		fmt.Printf("Failed to perform exec operation on instance: worker-%s-%s\n%v\n", goblin.ImageName, instanceNum, err)
		return err
	}

	if err = execOp.Wait(); err != nil {
		fmt.Printf("Error on exec Wait() operation: worker-%s-%s\n%v\n", goblin.ImageName, instanceNum, err)
		return err
	}
	return nil
}

func StopRequest(goblin *GoblinClient, instanceNum string) error {
	instanceName := fmt.Sprintf("worker-%s-%s", goblin.ImageName, instanceNum)
	stopReq := api.InstanceStatePut{
		Action:  "stop",
		Timeout: -1,
	}

	fmt.Printf("Stopping instance: worker-%s-%s\n", goblin.ImageName, instanceNum)

	stateOp, err := goblin.Server.UpdateInstanceState(instanceName, stopReq, "")
	if err != nil {
		fmt.Printf("Failed to perform stop operation on instance: worker-%s-%s\n%v\n", goblin.ImageName, instanceNum, err)
		return err
	}

	if err = stateOp.Wait(); err != nil {
		fmt.Printf("Error on stop Wait() operation: worker-%s-%s\n%v\n", goblin.ImageName, instanceNum, err)
		return err
	}
	return nil
}
func StopAllRequest(goblin *GoblinClient, instanceName string) error {
	stopReq := api.InstanceStatePut{
		Action:  "stop",
		Timeout: -1,
	}

	fmt.Printf("Stopping instance: %s\n", instanceName)

	stateOp, err := goblin.Server.UpdateInstanceState(instanceName, stopReq, "")
	if err != nil {
		fmt.Printf("Failed to perform stop operation on instance: worker-%s-%s\n%v\n", goblin.ImageName, instanceName, err)
		return err
	}

	if err = stateOp.Wait(); err != nil {
		fmt.Printf("Error on stop Wait() operation: %s\n", instanceName)
		return err
	}
	return nil
}

func StartRequest(goblin *GoblinClient, instanceNum string) error {
	instanceName := fmt.Sprintf("worker-%s-%s", goblin.ImageName, instanceNum)
	startReq := api.InstanceStatePut{
		Action:  "start",
		Timeout: -1,
	}

	fmt.Printf("Starting instance: worker-%s-%s\n", goblin.ImageName, instanceNum)

	stateOp, err := goblin.Server.UpdateInstanceState(instanceName, startReq, "")
	if err != nil {
		fmt.Printf("Failed to perform start operation on instance: worker-%s-%s\n%v\n", goblin.ImageName, instanceNum, err)
		return err
	}

	if err = stateOp.Wait(); err != nil {
		fmt.Printf("Error on start Wait() operation: worker-%s-%s\n%v\n", goblin.ImageName, instanceNum, err)
		return err
	}
	return nil
}

func DeleteRequest(goblin *GoblinClient, instanceNum string) error {
	instanceName := fmt.Sprintf("worker-%s-%s", goblin.ImageName, instanceNum)
	fmt.Printf("Deleting instance: worker-%s-%s\n", goblin.ImageName, instanceNum)

	deleteOp, err := goblin.Server.DeleteInstance(instanceName)
	if err != nil {
		fmt.Printf("Failed to perform delete operation on instance: worker-%s-%s\n%v\n", goblin.ImageName, instanceNum, err)
		return err
	}
	if err = deleteOp.Wait(); err != nil {
		fmt.Printf("Error on delete Wait() operation: worker-%s-%s\n%v\n", goblin.ImageName, instanceNum, err)
		return err
	}
	return nil
}

func DeleteAllRequest(goblin *GoblinClient, instanceName string) error {

	deleteOp, err := goblin.Server.DeleteInstance(instanceName)
	if err != nil {
		fmt.Printf("Failed to perform delete operation on instance: worker-%s-%s\n%v\n", goblin.ImageName, instanceName, err)
		return err
	}
	if err = deleteOp.Wait(); err != nil {
		fmt.Printf("Error on delete Wait() operation: worker-%s-%s\n%v\n", goblin.ImageName, instanceName, err)
		return err
	}
	return nil

}

func PushFileToWorker(goblin *GoblinClient, instanceNum string, file *os.File, fileName string) error {
	instanceName := fmt.Sprintf("worker-%s-%s", goblin.ImageName, instanceNum)
	containerDst := filepath.Join("/tmp", fileName)

	args := incus.InstanceFileArgs{
		Content:   file,
		GID:       0,
		UID:       0,
		Mode:      0644,
		Type:      "file",
		WriteMode: "overwrite",
	}

	if err := goblin.Server.CreateInstanceFile(instanceName, containerDst, args); err != nil {
		return fmt.Errorf("failed to push file to container %s: %w", instanceName, err)
	}

	return nil
}

func PullFileFromWorker(goblin *GoblinClient, instanceNum string, filePath string) error {
	instanceName := fmt.Sprintf("worker-%s-%s", goblin.ImageName, instanceNum)

	content, _, err := goblin.Server.GetInstanceFile(instanceName, filePath)
	if err != nil {
		return fmt.Errorf("Error getting instance file: %s: %w\n", instanceName, err)
	}
	defer content.Close()

	data, err := io.ReadAll(content)
	if err != nil {
		return fmt.Errorf("failed to read file content: %w", err)
	}
	fmt.Printf("\n--- Contents of %s ---\n", filePath)
	fmt.Println(string(data))
	fmt.Println("--------------------------------")

	return nil
}
