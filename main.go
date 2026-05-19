package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/oneto6/vcp-be/meet"
)

func main() {
	engin := gin.New()
	engin.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))
	engin.GET("/", func(c *gin.Context) {
		c.String(200, `Hello World`)
	})

	engin.GET("/meet", func(c *gin.Context) {
		meet := meet.GetMeet()
		c.JSON(200, meet)
	})
	engin.Run("127.0.0.1:0437")
}
