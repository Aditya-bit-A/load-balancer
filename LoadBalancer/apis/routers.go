package apis

import (
	"loadbalancer/server"

	"github.com/gin-gonic/gin"
)

func SetupRouter(manager *server.DefaultServerManager) *gin.Engine {
	r := gin.Default()
	r.GET("/rep", GetReplicasHandler(manager))
	r.POST("/add", AddServersHandler(manager))
	r.POST("/rem", RemoveServersHandler(manager))
	r.GET("/:path", RequestRedirectHandler(manager))
	return r
}

