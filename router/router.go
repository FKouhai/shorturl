// Package router has all the router configs and routes
package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"log"
	"net/http"
	"url_shortener/methods"
)

// InJSON struct to parse url
type InJSON struct {
	URL string `json:"url"`
}

// SetRouter initializes the gin router
func SetRouter() *gin.Engine {
	router := gin.Default()
	// refer to https://github.com/gin-gonic/gin/issues/2697#issuecomment-829071839
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"0.0.0.0/0"})
	router.RemoteIPHeaders = []string{"X-Forwarded-For", "X-Real-IP"}
	router.Static("/assets", "./assets")
	router.Use(gin.Recovery())
	router.Use(otelgin.Middleware("shorturl"))
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	router = HTMLRouter(router)
	router = APIGroups(router)

	return router
}

// HTMLRouter aaa
func HTMLRouter(router *gin.Engine) *gin.Engine {
	router.LoadHTMLGlob("views/*")
	return router
}

// APIGroups provides all the api groups for the router
func APIGroups(router *gin.Engine) *gin.Engine {
	api := router.Group("/api")
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	api.POST("/addRoute", func(ctx *gin.Context) {
		j := InJSON{}
		err := ctx.BindJSON(&j)
		println(j.URL)
		if err != nil {
			log.Println(err)
		}
		methods.CreateEntry(ctx, j.URL)
	})
	api.GET("/redirect/:path", func(ctx *gin.Context) {
		path := ctx.Params.ByName("path")
		methods.Redirect(ctx, path)
	})
	api.GET("/list", func(ctx *gin.Context) {
		v, err := methods.ListAll(ctx)
		if err != nil {
			log.Println(err)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"keys": fmt.Sprintf("%v", v),
		})

	})
	return router
}
