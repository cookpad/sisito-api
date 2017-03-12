package sisito

import (
	"gopkg.in/gin-gonic/gin.v1"
)

type Server struct {
	Engine *gin.Engine
	Router gin.IRouter
	Driver *Driver
}

func NewServer(config *Config, driver *Driver) (server *Server) {
	engine := gin.Default()

	server = &Server{
		Engine: engine,
		Router: engine,
		Driver: driver,
	}

	if len(config.User) > 0 {
		accounts := gin.Accounts{}

		for _, u := range config.User {
			accounts[u.Userid] = u.Password
		}

		server.Router = engine.Group("", gin.BasicAuth(accounts))
	}

	return
}

func (server *Server) Run() {
	router := server.Router

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	server.Engine.Run()
}
