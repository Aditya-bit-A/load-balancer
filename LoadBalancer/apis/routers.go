package apis

import (
	"loadbalancer/server"

	_ "loadbalancer/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(manager *server.DefaultServerManager) *gin.Engine {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/rep", GetReplicasHandler(manager))
	r.POST("/add", AddServersHandler(manager))
	r.POST("/rem", RemoveServersHandler(manager))
	r.GET("/:path", RequestRedirectHandler(manager))
	return r
}
