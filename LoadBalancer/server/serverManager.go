package server

import (
	"container/list"
	"fmt"
	"loadbalancer/docker"
	"loadbalancer/models"
	"loadbalancer/server/balancing_strategy"
	"sync"
)

// Server manager struct
type DefaultServerManager struct {
	serverMap        map[string]*list.Element
	serverList       *list.List
	loadBalancer     balancing_strategy.LoadBalancer
	mu               sync.RWMutex
	containerRuntime docker.ContainerRuntime
}

var (
	instance *DefaultServerManager
	once     sync.Once
)

// Private Constructor
func newDefaultServerManager() *DefaultServerManager {
	return &DefaultServerManager{
		serverMap:        make(map[string]*list.Element),
		serverList:       list.New(),
		containerRuntime: &docker.DockerContainerRuntime{},
	}
}

func GetManager() *DefaultServerManager {
	once.Do(func() {
		instance = newDefaultServerManager()
	})
	return instance
}

func (s *DefaultServerManager) SetLoadBalancer(lb balancing_strategy.LoadBalancer) {
	s.loadBalancer = lb
}

func (s *DefaultServerManager) SetContainerRuntime(cr docker.ContainerRuntime) {
	s.containerRuntime = cr
}

func (s *DefaultServerManager) insertNewInstanceOnSystem(hashValue int, inst models.ServerMetaData) {
	var currentElement *list.Element
	if s.serverList.Len() == 0 {
		currentElement = s.serverList.PushBack(inst)
	} else {
		inserted := false
		for e := s.serverList.Front(); e != nil; e = e.Next() {
			currentEHash := e.Value.(*models.ServerInstance).GetHashVal()

			if hashValue < currentEHash {
				currentElement = s.serverList.InsertBefore(inst, e)
				inserted = true
				break
			}
		}
		if !inserted {
			currentElement = s.serverList.PushBack(inst)
		}
	}

	if _, exists := s.serverMap[inst.GetServerID()]; !exists {
		//inst.guid = createNewProcess(inst)
		s.serverMap[inst.GetServerID()] = currentElement
	}
}

func (s *DefaultServerManager) AddNewServerInstances(payload models.AddServerReqPayload) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := 0; i < payload.N; i++ {
		ID := payload.HostNames[i]
		hashValue := models.SH(ID, 1)
		inst := NewServerInstance(ID, "Type1")
		// Insert this server Instance details on the load balancer (Can be any server instance)
		s.insertNewInstanceOnSystem(hashValue, inst) // Insert instance logic depends on the load balancer implementation
		s.containerRuntime.CreateNewContainer(inst)
	}
}

func (s *DefaultServerManager) SearchServerInstance(clientHash int) models.ServerMetaData {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.loadBalancer == nil {
		return nil
	}

	ctx := balancing_strategy.LoadBalancerContext{
		ServerMap:  s.serverMap,
		ServerList: s.serverList,
		ClientHash: clientHash,
	}
	server := s.loadBalancer.SelectServer(&ctx)
	return server
}

func (s *DefaultServerManager) RemoveServersFromList(payload models.AddServerReqPayload) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := 0; i < payload.N; i++ {
		ID := payload.HostNames[i]
		if elem, exists := s.serverMap[ID]; exists {
			s.serverList.Remove(elem)
			delete(s.serverMap, ID)
			s.containerRuntime.StopContainer(elem.Value.(*models.ServerInstance))
		}
	}
}

func (s *DefaultServerManager) ShowList() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for e := s.serverList.Front(); e != nil; e = e.Next() {
		fmt.Println(" -> ", e.Value.(*models.ServerInstance).ID)
	}
}

func (s *DefaultServerManager) GetCurrentServersList() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	replicas := make([]string, 0, len(s.serverMap))
	for _, val := range s.serverMap {
		replicas = append(replicas, val.Value.(*models.ServerInstance).ID)
	}
	//s.showList()
	return replicas
}
