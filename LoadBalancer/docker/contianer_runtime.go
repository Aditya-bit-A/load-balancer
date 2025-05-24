package docker

import (
	"bytes"
	"fmt"
	"loadbalancer/config"
	"loadbalancer/models"
	"log"
	"os/exec"
)

type ContainerRuntime interface {
	StopContainer(inst models.ServerMetaData) bool
	CreateNewContainer(inst models.ServerMetaData) string
}

// Define all the Container Runtime Structs
type DockerContainerRuntime struct{}

func (d *DockerContainerRuntime) StopContainer(inst models.ServerMetaData) bool {
	cmd := exec.Command("docker", "rm", "-f", inst.GetContainerID())
	err := cmd.Run()
	if err != nil {
		log.Println("Failed to stop container :", inst.GetServerID(), err)
		return false
	}
	return true
}

func (d *DockerContainerRuntime) CreateNewContainer(inst models.ServerMetaData) string {
	cmd := exec.Command(
		"docker", "run", "-d",
		"--network", config.GetEnv("NETWORK_NAME", "loadbalancer"),
		"--network-alias", inst.GetContainerHostName(),
		"--name", inst.GetContainerHostName(),
		"-p", "80:"+config.GetEnv("BACKEND_PORT", "8000"),
		"--env", "SERVER_ID="+inst.GetServerID(),
		"serverinst-app",
	)
	var out bytes.Buffer
	cmd.Stdout = &out
	//cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error creating new proccess:", err)
		return ""
	}

	return out.String()
}
