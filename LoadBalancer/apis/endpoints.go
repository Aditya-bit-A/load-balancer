package apis

import (
	"loadbalancer/models"
	"loadbalancer/server"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetReplicasHandler godoc
// @Summary      Get server replicas
// @Description  Returns the list of server replicas managed by the load balancer
// @Tags         servers
// @Produce      json
// @Success      200  {array}   models.ReplicasResponse
// @Failure      500  {object}  map[string]string
// @Router       /rep [get]
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

// AddServersHandler godoc
// @Summary      Add new server instances
// @Description  Adds new server instances to the load balancer
// @Tags         servers
// @Accept       json
// @Produce      json
// @Param        payload  body  models.AddServerReqPayload  true  "Add server payload"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /add [post]
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

// RemoveServersHandler godoc
// @Summary      Remove server instances
// @Description  Removes one or more server instances from the load balancer
// @Tags         servers
// @Accept       json
// @Produce      json
// @Param        payload  body  models.AddServerReqPayload  true  "Payload to remove server instances"
// @Success      200  {object}  models.ReplicasResponse
// @Failure      400  {object}  map[string]string
// @Router       /rem [post]
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

// RequestRedirectHandler godoc
// @Summary      Redirect request to a backend server
// @Description  Forwards the client request to a selected server based on load balancing
// @Tags         routing
// @Accept       json
// @Produce      json
// @Param        path  path  string  true  "Path to be forwarded"
// @Success      200  {object}  string
// @Failure      404  {object}  map[string]string
// @Router       /{path} [get]
func RequestRedirectHandler(manager *server.DefaultServerManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the path parameter as a string
		fwdPath := c.Param("path")

		// Check if such a path exists in server endpoints
		requestId := models.GenerateRequestId()
		clientHash := models.H(requestId)

		serverInst := manager.SearchServerInstance(clientHash)
		log.Printf("Reuqest Redirect :- Request ID: %s, Client Hash: %d, Selected Server: %s", requestId, clientHash, serverInst.GetContainerHostName())
		if serverInst == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No suitable server found for redirection"})
			return
		}
		// forward request to the server
		forwardRequestToServer(c, fwdPath, serverInst)
	}
}
