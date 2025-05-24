package main

import (
	"loadbalancer/apis"
	"loadbalancer/docker"
	"loadbalancer/server"
	"loadbalancer/server/balancing_strategy"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment")
	}
	serverManager := server.GetManager()
	serverManager.SetLoadBalancer(&balancing_strategy.HashingLoadBalancer{})
	serverManager.SetContainerRuntime(&docker.DockerContainerRuntime{})
	
	router := apis.SetupRouter(serverManager)
	router.Run(":8080")

}
