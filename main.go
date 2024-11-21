package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Monitoring, logging, -> apakah proses berhasil, response seperti apa
// analisis, apakah prosesnya cepat, kalau lama itu di codingan mana, dst

// analitic, google analitic

// Auth

// rate limiter, user ini batasi max request 10 req per detik

// logic bisnis

// membuat logic lebih efisien dengan middleware
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		timeNow := time.Now()

		c.Set("example", "123456")

		// forward / melanjutkan ke handler, middleware selanjutnya
		c.Next()
		log.Println("Time elapsed: ", time.Since(timeNow))
	}
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ngirim query url dengan key berupa token, value
		// request
		token := c.Query("token")

		if token == "" {
			// batalkan proses
			c.AbortWithStatusJSON(401, gin.H{"message": "unauthorized"})
			return
		}

		c.Next()

		// response
	}
}

// key value
// init middleware
// lempar data dari middleware ke handler
// lempar middleware ke middleware lain

func main() {
	r := gin.Default()

	// set middleware global
	// r.Use(Logger())

	// set middleware di group
	v1 := r.Group("/v1", Logger())
	{
		v1.GET("/example", func(c *gin.Context) {
			example := c.MustGet("example").(string)

			c.JSON(http.StatusOK, gin.H{"example": example})
		})

		v1.GET("/hello", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "hello"})
		})

		v1.GET("/index", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "index"})
		})
	}

	v1Private := r.Group("/v1/private", Logger(), Auth())
	{
		v1Private.GET("/home", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "home"})
		})
	}

	//

	r.Run(":8080")
}
