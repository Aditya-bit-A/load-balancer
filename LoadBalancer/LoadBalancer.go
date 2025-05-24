package main

import (
	"bytes"
	"container/list"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

const K int = 512 // SLOTS for Servers
const VSERV int = 9
const NETWORK_NAME string = "lb_network"

var (
	serversMap     = make(map[string]*list.Element)
	serversList    = list.New()
	serverMapMutex sync.RWMutex
)

func H(i int) int {
	hash := (i*i + 2*i + 17) % K
	return hash
}
func SH(i int, j int) int {
	serverHash := (i*i + j*j + 2*j + 25) % K
	return serverHash
}

type serverInstance struct {
	ID      int
	HashVal int
	guid    string
}

func (s *serverInstance) getContainerId() string {
	return s.guid
}
func (s *serverInstance) getContainerHostName() string {
	return "serverinst" + strconv.Itoa(s.ID)
}
func (s *serverInstance) getServerID() string {
	return "Server_" + strconv.Itoa(s.ID)
}

type ReplicasResponse struct {
	Message struct {
		N        int   `json:"N"`
		Replicas []int `json:"replicas"`
	} `json:"message"`
	Status string `json:"status"`
}

type AddServerReqPayload struct {
	N         int   `json:"n"`
	HostNames []int `json:"hostnames"`
}

func forwardRequestToServer(c *gin.Context, fwdPath string, serverInst *serverInstance) {
	// Construct backend URL
	backendURL := "http://" + serverInst.getContainerHostName() + ":8000/" + fwdPath

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

func showList() {
	for e := serversList.Front(); e != nil; e = e.Next() {
		fmt.Println(" -> ", e.Value.(serverInstance).ID)
	}
}
func getCurrentServersList() []int {
	replicas := make([]int, 0, K)
	for _, val := range serversMap {
		replicas = append(replicas, val.Value.(serverInstance).ID)
	}
	showList()
	return replicas
}

func stopContainer(inst *serverInstance) {
	cmd := exec.Command("docker", "rm", "-f", inst.getContainerId())
	err := cmd.Run()
	if err != nil {
		log.Println("Failed to stop container :", inst.getServerID(), err)
	}
}

func createNewProcess(inst *serverInstance) string {
	cmd := exec.Command(
		"docker", "run", "-d",
		"--network", NETWORK_NAME,
		"--network-alias", inst.getContainerHostName(),
		"--name", inst.getContainerHostName(),
		"-p", "80:80",
		"--env", "SERVER_ID="+inst.getServerID(),
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
func addNewServerToList(payload AddServerReqPayload) {

	for i := 0; i < payload.N; i++ {
		var ID = payload.HostNames[i]
		var currentElement *list.Element = nil
		hashValue := SH(payload.HostNames[i], 1)
		fmt.Println("Payload:", payload.HostNames[i])
		if serversList.Len() == 0 {
			currentElement = serversList.PushBack(&serverInstance{ID: payload.HostNames[i], HashVal: hashValue})
		} else {
			for e := serversList.Front(); e != nil; e = e.Next() {
				currentEHash := e.Value.(*serverInstance).HashVal

				if hashValue > currentEHash {
					currentElement = serversList.InsertAfter(&serverInstance{ID: payload.HostNames[i], HashVal: hashValue}, e)
					fmt.Println("Added After:", currentElement.Value.(*serverInstance).ID)
					break
				}
			}
		}
		inst := currentElement.Value.(*serverInstance)
		_, full := serversMap[inst.getServerID()]
		container_id := createNewProcess(currentElement.Value.(*serverInstance))
		// handle error here
		inst.guid = container_id
		if !full {
			serversMap[inst.getServerID()] = currentElement
		}

	}

}
func searchServerInstance(clientHash int) *serverInstance {
	if serversList.Len() == 0 {
		return nil
	}
	for e := serversList.Front(); e != nil; e = e.Next() {
		server := e.Value.(*serverInstance)
		currentEHash := server.HashVal
		serverID := server.ID

		fmt.Println("Server Hash:", currentEHash)

		if clientHash < currentEHash {
			response := fmt.Sprintf("Request Redirected to : %d", serverID)
			fmt.Println(response)
			// forward request to the server
			return server
		}
	}
	return serversList.Front().Value.(*serverInstance)
}
func removeServersFromList(payload AddServerReqPayload) {

	for i := 0; i < payload.N; i++ {
		var ID = payload.HostNames[i]
		serversList.Remove(serversMap[ID])
		delete(serversMap, ID)
	}
}

func getReplicas(c *gin.Context) {

	serverMapMutex.RLock()
	replicas := getCurrentServersList()
	serverMapMutex.RUnlock()

	response := ReplicasResponse{
		Status: "successful",
	}
	response.Message.N = len(replicas)
	response.Message.Replicas = replicas
}

func getReplicasGin(c *gin.Context) {
	// Get the current list of servers
	replicas := getCurrentServersList()

	// Build the response
	response := ReplicasResponse{
		Status: "successful",
	}
	response.Message.N = len(replicas)
	response.Message.Replicas = replicas

	// Send the JSON response
	c.JSON(http.StatusOK, response)
}

func addServersGin(c *gin.Context) {
	var payload AddServerReqPayload

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

	serverMapMutex.Lock()
	// Perform the server addition
	addNewServerToList(payload)
	serverMapMutex.Unlock()

	// Get updated list of replicas
	serverMapMutex.RLock()
	replicas := getCurrentServersList()
	serverMapMutex.RUnlock()

	// Build and return the response
	response := ReplicasResponse{
		Status: "successful",
	}
	response.Message.N = len(replicas)
	response.Message.Replicas = replicas

	c.JSON(http.StatusOK, response)
}

func removeServersGin(c *gin.Context) {
	var payload AddServerReqPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body format."})
		return
	}

	if payload.N != len(payload.HostNames) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Length not matching with list."})
		return
	}

	serverMapMutex.Lock()
	removeServersFromList(payload)
	serverMapMutex.Unlock()

	replicas := []int{1} // Dummy replicas
	response := ReplicasResponse{
		Status: "successful",
	}
	response.Message.N = len(replicas)
	response.Message.Replicas = replicas

	c.JSON(http.StatusOK, response)
}

func requestRedirectGin(c *gin.Context) {
	// Get the path parameter as a string
	fwdPath := c.Param("path")

	// Check if such a path exists in server endpoints

	clientHash := H(requestId)
	//fmt.Println("Request Hash:", clientHash)

	serverInst := searchServerInstance(clientHash)
	if serverInst == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No suitable server found for redirection"})
		return
	}
	// forward request to the server
	var resp = forwardRequestToServer(c, fwdPath, serverInst)
	c.JSON(http.StatusOK, resp)
}

func main() {
	router := gin.Default()

	//serversMap = make(map[string]*list.Element)
	router.GET("/rep", getReplicasGin)
	router.POST("/add", addServersGin)
	router.POST("/rem", removeServersGin)
	router.GET("/:path", requestRedirectGin)

	var initialPayload AddServerReqPayload
	initialPayload.N = 3
	initialPayload.HostNames = make([]int, 0, 3)
	for i := 0; i < 3; i++ {
		initialPayload.HostNames = append(initialPayload.HostNames, i)
	}
	addNewServerToList(initialPayload)
	log.Printf("Added 3 servers initially.. ")
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000" // default port
	}

	log.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
