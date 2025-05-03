// Package main contains the api groups etc...
package main

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
	_ "modernc.org/sqlite"
	"url_shortener/router"
	"url_shortener/tracer"
)

//go:embed schemas/schema.sql
var ddl string

// Initializes the database connection and executes DDL statements.
func init() {
	ctx := context.Background()

	db, err := sql.Open("sqlite", "sqlitedb")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	err = db.PingContext(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.ExecContext(ctx, ddl)
	if err != nil {
		log.Println(err)
		return
	}
}
func init() {
	ctx := context.Background()

	db, err := sql.Open("sqlite", "sqlitedb")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	err = db.PingContext(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.ExecContext(ctx, ddl)
	if err != nil {
		log.Println(err)
		return
	}
}

// InJSON public struct

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
	router := router.SetRouter()
	err = router.Run("0.0.0.0:7880")
	if err != nil {
		log.Println(err)
	}

}
