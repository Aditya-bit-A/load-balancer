package models

type AddServerReqPayload struct {
	N         int      `json:"n"`
	HostNames []string `json:"hostnames"`
}

type ReplicasResponse struct {
	Message struct {
		N        int      `json:"N"`
		Replicas []string `json:"replicas"`
	} `json:"message"`
	Status string `json:"status"`
}
