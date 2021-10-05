package WebServer

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type WebServer struct {
	Port int
}

func NewWebServer(port int) *WebServer {
	server := &WebServer{
		Port: port,
	}
	return server
}

func (webServer *WebServer) Start() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(fmt.Sprintf(":%d", webServer.Port))
}
