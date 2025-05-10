package apis

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func addServers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload AddServerReqPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request body format.", http.StatusBadRequest)
		return
	}

	if payload.N != len(payload.HostNames) {
		http.Error(w, "Length not matching with list.", http.StatusBadRequest)
	}

	log.Printf("Adding %d new servers: %v\n", payload.N, payload.HostNames)
	addNewServerToList(payload)

	w.Header().Set("Content-Type", "application./json")
	w.WriteHeader(http.StatusOK)

	replicas := getCurrentServersList()

	response := ReplicasResponse{
		Status: "successful",
	}
	response.Message.N = len(replicas)
	response.Message.Replicas = replicas

}

func removeServers(w http.ResponseWriter, r *http.Request) {
	var payload AddServerReqPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request body format.", http.StatusBadRequest)
		return
	}

	if payload.N != len(payload.HostNames) {
		http.Error(w, "Length not matching with list.", http.StatusBadRequest)
	}

	removeServersFromList(payload)

	w.Header().Set("Content-Type", "application./json")
	w.WriteHeader(http.StatusOK)
	replicas := []int{1}
	response := ReplicasResponse{
		Status: "successful",
	}
	response.Message.N = len(replicas)
	response.Message.Replicas = replicas

}
func requestRedirect(w http.ResponseWriter, r *http.Request) {
	// Trim the leading and trailing slashes
	path := strings.Trim(r.URL.Path, "/")

	if path == "" {
		fmt.Fprintln(w, "Welcome to the root!")
		return
	}
	numpath, _ := strconv.Atoi(path)
	clientHash := H(numpath)
	keys := make([]int, 0, len(serversMap))

	for k, _ := range serversMap {
		keys = append(keys, k)
	}
	fmt.Println("Req ", clientHash)
	for e := serversList.Front(); e != nil; e = e.Next() {
		currentEHash := e.Value.(serverInstance).HashVal
		var Id = e.Value.(serverInstance).ID
		fmt.Println(currentEHash)
		if clientHash < currentEHash {
			response := fmt.Sprintf("Request Redirected to : %d", Id)
			fmt.Println(response)
			w.Write([]byte(response))
			return
		}
	}

	// Handle dynamic part

}
