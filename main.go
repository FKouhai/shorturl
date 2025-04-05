// Package main contains the api groups etc...
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"url_shortener/methods"

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
	db.Close()
	//_ = dbpkg.New(db)

}

// InJSON public struct
type InJSON struct {
	Url string `json:"url"`
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
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"keys": fmt.Sprintf("%v", v),
		})

	})

	router.Run("0.0.0.0:7880")

}
