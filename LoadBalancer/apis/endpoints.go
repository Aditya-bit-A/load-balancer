package apis

import (
	"loadbalancer/models"
	"loadbalancer/server"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetReplicasHandler(manager *server.DefaultServerManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the current list of servers
		replicas := manager.GetCurrentServersList()

		// Build the response
		response := models.ReplicasResponse{
			Status: "successful",
		}
		response.Message.N = len(replicas)
		response.Message.Replicas = replicas

		// Send the JSON response
		c.JSON(http.StatusOK, response)
	}
}

func AddServersHandler(manager *server.DefaultServerManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload models.AddServerReqPayload

		// Bind JSON body to payload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body format."})
			return
		}

		// Validate input
		if payload.N != len(payload.HostNames) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Length not matching with list."})
			return
		}

		// Log the operation
		log.Printf("Adding %d new servers: %v\n", payload.N, payload.HostNames)

		// Perform the server addition
		manager.AddNewServerInstances(payload)

		// Get updated list of replicas
		replicas := manager.GetCurrentServersList()

		// Build and return the response
		response := models.ReplicasResponse{
			Status: "successful",
		}
		response.Message.N = len(replicas)
		response.Message.Replicas = replicas

		c.JSON(http.StatusOK, response)
	}
}

func RemoveServersHandler(manager *server.DefaultServerManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload models.AddServerReqPayload

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body format."})
			return
		}

		if payload.N != len(payload.HostNames) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Length not matching with list."})
			return
		}

		manager.RemoveServersFromList(payload)

		replicas := []string{"1"} // Dummy replicas
		response := models.ReplicasResponse{
			Status: "successful",
		}
		response.Message.N = len(replicas)
		response.Message.Replicas = replicas

		c.JSON(http.StatusOK, response)
	}
}

func RequestRedirectHandler(manager *server.DefaultServerManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the path parameter as a string
		fwdPath := c.Param("path")

		// Check if such a path exists in server endpoints
		requestId := models.GenerateRequestId()
		clientHash := models.H(requestId)

		serverInst := manager.SearchServerInstance(clientHash)
		if serverInst == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No suitable server found for redirection"})
			return
		}
		// forward request to the server
		forwardRequestToServer(c, fwdPath, serverInst)
	}
}
