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

func (server *Server) ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (server *Server) bounced(c *gin.Context) {
	recipient := c.Query("recipient")
	digest := c.Query("digest")
	senderdomain := c.Query("senderdomain")

	if recipient == "" && digest == "" {
		c.JSON(400, gin.H{
			"message": `"recipient" or "digest" is not present`,
		})
	} else if recipient == "" && digest == "" {
		c.JSON(400, gin.H{
			"message": ` Cannot pass both "recipient" and "digest"`,
		})
	} else {
		var name string
		var value string

		if recipient != "" && digest == "" {
			name = "recipient"
			value = recipient
		} else if recipient == "" || digest != "" {
			name = "digest"
			value = digest
		}

		bounced, err := server.Driver.IsBounced(name, value, senderdomain)

		if err != nil {
			panic(err)
		}

		c.JSON(200, gin.H{
			"bounced": bounced,
		})
	}
}

func (server *Server) Run() {
	engine := server.Engine
	router := server.Router

	engine.GET("/ping", server.ping)
	router.GET("/bounced", server.bounced)

	engine.Run()
}
