// Package main contains the api groups etc...
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"url_shortener/methods"
	"url_shortener/tracer"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	_ "modernc.org/sqlite"

	"context"
	_ "embed"

	"github.com/gin-gonic/gin"
)

//go:embed schemas/schema.sql
var ddl string

func init() {
	ctx := context.Background()
	db, err := sql.Open("sqlite", "sqlitedb")

	if err != nil {
		log.Fatalln(err)
	}

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		log.Println(err)
	}
	err = db.Close()
	if err != nil {
		log.Println(err)
	}

}

// InJSON public struct
type InJSON struct {
	Url string `json:"url"`
}

func main() {
	tp, err := tracer.InitTracer()
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(otelgin.Middleware("shorturl"))
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	router.POST("/addRoute", func(ctx *gin.Context) {
		j := InJSON{}
		err := ctx.BindJSON(&j)
		println(j.Url)
		if err != nil {
			log.Println(err)
		}
		methods.CreateEntry(ctx, j.Url)
	})

	router.GET("/redirect/:path", func(ctx *gin.Context) {
		path := ctx.Params.ByName("path")
		methods.Redirect(ctx, path)
	})
	router.GET("/list", func(ctx *gin.Context) {
		v, err := methods.ListAll()
		if err != nil {
			log.Println(err)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"keys": fmt.Sprintf("%v", v),
		})

	})

	err = router.Run("0.0.0.0:7880")
	if err != nil {
		log.Println(err)
	}

}
