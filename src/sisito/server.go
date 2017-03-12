package sisito

import (
	"gopkg.in/gin-gonic/gin.v1"
)

type Server struct {
	Engine *gin.Engine
	Driver *driver
}

func NewServer(driver *Driver) (server *Server) {
	server = &Server{
		Engine: gin.Default(),
		Driver: driver,
	}
}

func (server *Server) Run() {
	engine := server.Engine

	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	engine.Run()
}
