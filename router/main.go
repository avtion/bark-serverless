package router

import "github.com/gin-gonic/gin"

var routerQueue = make([]func(r *gin.Engine), 0)

func AppendRouter(fns ...func(r *gin.Engine)) {
	routerQueue = append(routerQueue, fns...)
}

func SetupRouter() *gin.Engine {
	server := gin.Default()
	for _, fn := range routerQueue {
		fn(server)
	}
	return server
}
