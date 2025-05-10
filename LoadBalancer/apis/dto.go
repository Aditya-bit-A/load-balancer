package apis

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
