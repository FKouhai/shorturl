// Package methods provides all the function helpers for gin handlers
package methods

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"url_shortener/db"
	"url_shortener/mem_storage"

	_ "modernc.org/sqlite"

	"github.com/gin-gonic/gin"
)

func open() (*db.Queries, *sql.DB) {
	d, err := sql.Open("sqlite", "sqlitedb")
	if err != nil {
		log.Println(err)
	}
	q := db.New(d)
	return q, d

}

func lookUPId(url string) (string, error) {
	query, _ := open()
	ctx := context.Background()
	key, err := query.GetUrlId(ctx, url)
	println(key)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(key), nil

}

// ListAll method returns a json encoded blob with all the created URL's
func ListAll() (any, error) {
	query, _ := open()
	ctx := context.Background()
	key, err := query.GetUrls(ctx)
	if err != nil {
		return nil, err
	}
	return key, nil

}

// CreateEntry adds a new entry in the database based on the api path given
func CreateEntry(c *gin.Context, url string) {
	v, err := lookUPId(url)
	if err != nil {
		log.Println(err)
	}
	if !isCached(v) || !isLTS(url) {
		log.Println("already created")
	}
	query, _ := open()
	ctx := context.Background()
	log.Println(url)
	_, err = query.CreateUrl(ctx, url)

	if err != nil {
		log.Println(err)
	}

	err = memstorage.SetValue(v, url)
	if err != nil {
		log.Println(err)
	}

	c.Done()
}

func isLTS(url string) bool {
	query, _ := open()
	ctx := context.Background()
	_, err := query.GetUrlId(ctx, url)
	return err == nil

}

// isCached checks against valkey if a given string is present in it's memory
func isCached(url string) bool {
	_, err := memstorage.GetKey(url)
	return err == nil

}
func toInt64(v string) int64 {
	var i int64
	_, err := fmt.Sscan(v, &i)
	if err != nil {
		return -1
	}
	return i
}

// Redirect is used to forward the traffic on the given path to the actual service
func Redirect(c *gin.Context, path string) {
	var url string
	q, _ := open()
	ctx := context.Background()
	if !isCached(path) {
		i := toInt64(path)
		d, _ := q.GetUrlData(ctx, i)
		url = d.Name
		id := string(d.ID)
		log.Println("path is not cached, performing sql query to get the dst addr")
		err := memstorage.SetValue(id, url)
		if err != nil {
			log.Fatalln("Unable to create entry in valkey ->", err)
		}
		c.Redirect(http.StatusPermanentRedirect, url)
	} else {
		url, err := memstorage.GetKey(path)
		if err != nil {
			log.Println(err)
		}
		log.Println("url has been cached")
		c.Redirect(http.StatusPermanentRedirect, url)
	}

}
