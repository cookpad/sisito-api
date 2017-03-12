package sisito

import (
	"gopkg.in/gin-gonic/gin.v1"
	"log"
)

func Debugf(format string, args ...interface{}) {
	if gin.Mode() == "debug" {
		log.Printf(format, args...)
	}
}
