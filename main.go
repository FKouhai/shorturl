// Package main contains the api groups etc...
package main

import (
	"database/sql"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"url_shortener/methods"

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
	db.Close()
	//_ = dbpkg.New(db)

}

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

	router.GET("/redirect/:path", func(ctx *gin.Context) {
		path := ctx.Params.ByName("path")
		methods.Redirect(ctx, path)
	})

	router.Run("0.0.0.0:7880")

}
