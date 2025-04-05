// Package main contains the api groups etc...
package main

import (
	"net/http"
	"url_shortener/methods"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(gin.Recovery())
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})
	router.POST("/addRoute", func(ctx *gin.Context) {
		methods.CreateEntry(ctx, "https://google.com")
	})
	router.GET("/redirect", func(ctx *gin.Context) {
		methods.Redirect(ctx, "0")

	})
	router.Run("0.0.0.0:7880")

}
