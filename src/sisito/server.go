package sisito

import (
	"gopkg.in/gin-gonic/gin.v1"
	"strconv"
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

func (server *Server) recent(c *gin.Context) {
	recipient := c.Query("recipient")
	digest := c.Query("digest")
	senderdomain := c.Query("senderdomain")

	if recipient == "" && digest == "" {
		c.JSON(400, gin.H{
			"message": `"recipient" or "digest" is not present`,
		})
	} else if recipient == "" && digest == "" {
		c.JSON(400, gin.H{
			"message": `Cannot pass both "recipient" and "digest"`,
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

		bounced, err := server.Driver.recentlyBounced(name, value, senderdomain)

		if err != nil {
			panic(err)
		}

		if len(bounced) > 0 {
			row := bounced[0]

			softbounce := false

			if row.Softbounce == 1 {
				softbounce = true
			}

			whitelisted := false

			if row.Whitelisted == 1 {
				whitelisted = true
			}

			c.JSON(200, gin.H{
				"timestamp":      row.Timestamp,
				"lhost":          row.Lhost,
				"rhost":          row.Rhost,
				"alias":          row.Alias,
				"reason":         row.Reason,
				"subject":        row.Subject,
				"messageid":      row.Messageid,
				"smtpagent":      row.Smtpagent,
				"softbounce":     softbounce,
				"smtpcommand":    row.Smtpcommand,
				"destination":    row.Destination,
				"senderdomain":   row.Senderdomain,
				"diagnosticcode": row.Diagnosticcode,
				"deliverystatus": row.Deliverystatus,
				"timezoneoffset": row.Timezoneoffset,
				"addresser":      row.Addresser,
				"recipient":      row.Recipient,
				"digest":         row.Digest,
				"created_at":     row.CreatedAt,
				"updated_at":     row.UpdatedAt,
				"whitelisted":    whitelisted,
			})
		} else {
			c.JSON(204, gin.H{})
		}
	}
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
			"message": `Cannot pass both "recipient" and "digest"`,
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

		bounced, err := server.Driver.isBounced(name, value, senderdomain)

		if err != nil {
			panic(err)
		}

		c.JSON(200, gin.H{
			"bounced": bounced,
		})
	}
}

func (server *Server) blacklist(c *gin.Context) {
	var limit uint64 = 0
	var err error

	senderdomain := c.Query("senderdomain")
	limitStr := c.Query("limit")

	if limitStr != "" {
		limit, err = strconv.ParseUint(limitStr, 10, 64)

		if err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})

			return
		}
	}

	var recipients []string
	recipients, err = server.Driver.blacklistRecipients(senderdomain, limit)

	if err != nil {
		panic(err)
	}

	c.JSON(200, gin.H{
		"recipients": recipients,
	})
}

func (server *Server) Run() {
	engine := server.Engine
	router := server.Router

	engine.GET("/ping", server.ping)
	router.GET("/recent", server.recent)
	router.GET("/bounced", server.bounced)
	router.GET("/blacklist", server.blacklist)

	engine.Run()
}
