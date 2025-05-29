// @title           Load Balancer API
// @version         1.0
// @description     This is the API documentation for the load balancer
// @BasePath        /
package main

import (
	"loadbalancer/apis"
	"loadbalancer/config"
	"loadbalancer/docker"
	"loadbalancer/server"
	"loadbalancer/server/balancing_strategy"
	"log"

	_ "loadbalancer/docs" // Import the generated docs package

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
	serverManager.InitalizeServers(3)
	router := apis.SetupRouter(serverManager)
	router.Run(":" + config.GetEnv("PORT", "5000")) // Default port is 8080 if not specified in .env

}
