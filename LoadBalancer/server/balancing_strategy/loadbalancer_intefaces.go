package balancing_strategy

import (
	"container/list"
	"loadbalancer/models"
	"log"
)

// interfaces/loadbalancer.go
type LoadBalancerContext struct {
	ServerMap  map[string]*list.Element
	ServerList *list.List
	ClientHash int
}

type LoadBalancer interface {
	SelectServer(ctx *LoadBalancerContext) models.ServerMetaData
}

// All the Load Balancer Strategy
type HashingLoadBalancer struct{}

func (lb *HashingLoadBalancer) SelectServer(ctx *LoadBalancerContext) models.ServerMetaData {
	if ctx.ServerList.Len() == 0 {
		return nil
	}
	for e := ctx.ServerList.Front(); e != nil; e = e.Next() {
		server := e.Value.(*models.ServerInstance)
		if ctx.ClientHash < server.HashVal {
			log.Printf("Request Redirected to : %s\n", server.ID)
			return server
		}
	}
	return ctx.ServerList.Front().Value.(*models.ServerInstance)
}
