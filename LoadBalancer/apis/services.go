package apis

import (
	"io"
	"loadbalancer/config"
	"loadbalancer/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func forwardRequestToServer(c *gin.Context, fwdPath string, serverInst models.ServerMetaData) {
	// Construct backend URL
	backendURL := "http://" + serverInst.GetContainerHostName() + ":" + config.GetEnv("BACKEND_PORT", "80") + "/" + fwdPath

	// Create new request with context from the client request
	req, err := http.NewRequestWithContext(c.Request.Context(), c.Request.Method, backendURL, c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to create request to backend: %v", err)
		return
	}

	// Copy headers
	for k, v := range c.Request.Header {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}

	// Send request to backend
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.String(http.StatusBadGateway, "Error contacting backend server: %v", err)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for k, v := range resp.Header {
		for _, vv := range v {
			c.Writer.Header().Add(k, vv)
		}
	}

	c.Status(resp.StatusCode)
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		log.Printf("Error copying response body: %v", err)
	}
}
