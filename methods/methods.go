// Package methods provides all the function helpers for gin handlers
package methods

import (
	"log"
	"net/http"
	"url_shortener/mem_storage"

	"github.com/gin-gonic/gin"
)

// CreateEntry adds a new entry in the database based on the api path given
func CreateEntry(c *gin.Context, url string) {
	if !isCached("0") {
		log.Println("already cached")
		return
	}

	err := memstorage.SetValue("0", url)
	if err != nil {
		log.Println(err)
	}

	c.Done()
}

// isCached checks against valkey if a given string is present in it's memory
func isCached(url string) bool {
	_, err := memstorage.GetKey(url)

	if err != nil {
		return false
	}

	return true
}

// Redirect is used to forward the traffic on the given path to the actual service
func Redirect(c *gin.Context, path string) {
	if !isCached(path) {
		// perform a sql query to fetch the url to redirect to
		log.Println("path is not cached, performing sql query to get the dst addr")
	}

	url, err := memstorage.GetKey(path)

	if err != nil {
		log.Println(err)
		return
	}

	c.Redirect(http.StatusPermanentRedirect, url)
}
