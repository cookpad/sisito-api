package sisito

import (
	"github.com/fvbock/endless"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"strconv"
)

type Server struct {
	Engine *gin.Engine
	Driver *Driver
}

func NewServer(config *Config, driver *Driver, out io.Writer) (server *Server) {
	engine := gin.New()
	engine.Use(gin.Recovery())

	if config.Server.Gzip {
		engine.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	logger := gin.LoggerWithWriter(out)
	engine.Use(logger)

	server = &Server{
		Engine: engine,
		Driver: driver,
	}

	var router, router_with_auth gin.IRouter
	router = engine.Group(config.Server.Prefix)
	router_with_auth = router

	if len(config.User) > 0 {
		accounts := gin.Accounts{}

		for _, u := range config.User {
			accounts[u.Userid] = u.Password
		}

		router_with_auth = engine.Group(config.Server.Prefix, gin.BasicAuth(accounts))
	}

	server.Engine.GET("/ping", server.Ping)

	if config.Authz.Recent {
		router_with_auth.GET("/recent", server.Recent)
	} else {
		router.GET("/recent", server.Recent)
	}

	if config.Authz.Listed {
		router_with_auth.GET("/listed", server.Listed)
	} else {
		router.GET("/listed", server.Listed)
	}

	if config.Authz.Blacklist {
		router_with_auth.GET("/blacklist", server.Blacklist)
	} else {
		router.GET("/blacklist", server.Blacklist)
	}

	return
}

func (server *Server) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (server *Server) Recent(c *gin.Context) {
	recipient := c.Query("recipient")
	digest := c.Query("digest")
	senderdomain := c.Query("senderdomain")
	filterStr := c.Query("filter")

	if recipient == "" && digest == "" {
		c.JSON(400, gin.H{
			"message": `"recipient" or "digest" is not present`,
		})
	} else if recipient != "" && digest != "" {
		c.JSON(400, gin.H{
			"message": `Cannot pass both "recipient" and "digest"`,
		})
	} else {
		var name string
		var value string
		var useFilter bool

		if recipient != "" && digest == "" {
			name = "recipient"
			value = recipient
		} else if recipient == "" || digest != "" {
			name = "digest"
			value = digest
		}

		if filterStr != "" {
			var err error
			useFilter, err = strconv.ParseBool(filterStr)

			if err != nil {
				c.JSON(400, gin.H{
					"message": err.Error(),
				})

				return
			}
		} else {
			useFilter = true
		}

		listed, err := server.Driver.RecentlyListed(name, value, senderdomain, useFilter)

		if err != nil {
			panic(err)
		}

		if len(listed) > 0 {
			row := listed[0]

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
			c.JSON(200, gin.H{})
		}
	}
}

func (server *Server) Listed(c *gin.Context) {
	recipient := c.Query("recipient")
	digest := c.Query("digest")
	senderdomain := c.Query("senderdomain")
	filterStr := c.Query("filter")

	if recipient == "" && digest == "" {
		c.JSON(400, gin.H{
			"message": `"recipient" or "digest" is not present`,
		})
	} else if recipient != "" && digest != "" {
		c.JSON(400, gin.H{
			"message": `Cannot pass both "recipient" and "digest"`,
		})
	} else {
		var name string
		var value string
		var useFilter bool

		if recipient != "" && digest == "" {
			name = "recipient"
			value = recipient
		} else if recipient == "" || digest != "" {
			name = "digest"
			value = digest
		}

		if filterStr != "" {
			var err error
			useFilter, err = strconv.ParseBool(filterStr)

			if err != nil {
				c.JSON(400, gin.H{
					"message": err.Error(),
				})

				return
			}
		} else {
			useFilter = true
		}

		listed, err := server.Driver.Listed(name, value, senderdomain, useFilter)

		if err != nil {
			panic(err)
		}

		c.JSON(200, gin.H{
			"listed": listed,
		})
	}
}

func (server *Server) Blacklist(c *gin.Context) {
	var softbounce *bool
	var limit, offset uint64
	var useFilter bool
	var err error

	senderdomain := c.Query("senderdomain")
	reasons := c.QueryArray("reason")
	softbounceStr := c.Query("softbounce")
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")
	filterStr := c.Query("filter")

	if softbounceStr != "" {
		softbounce = new(bool)
		*softbounce, err = strconv.ParseBool(softbounceStr)

		if err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})

			return
		}
	}

	if limitStr != "" {
		limit, err = strconv.ParseUint(limitStr, 10, 64)

		if err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})

			return
		}
	}

	if offsetStr != "" {
		offset, err = strconv.ParseUint(offsetStr, 10, 64)

		if err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})

			return
		}
	}

	if filterStr != "" {
		useFilter, err = strconv.ParseBool(filterStr)

		if err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})

			return
		}
	} else {
		useFilter = true
	}

	var recipients []string
	recipients, err = server.Driver.BlacklistRecipients(senderdomain, reasons, softbounce, limit, offset, useFilter)

	if err != nil {
		panic(err)
	}

	c.JSON(200, gin.H{
		"recipients": recipients,
	})
}

func (server *Server) Run() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	endless.ListenAndServe(":"+port, server.Engine)
}
