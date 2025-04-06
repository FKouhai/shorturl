// Package methods provides all the function helpers for gin handlers
package methods

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"url_shortener/db"
	"url_shortener/mem_storage"
	"url_shortener/tracer"

	_ "modernc.org/sqlite"

	"github.com/gin-gonic/gin"
)

func open(ctx context.Context) (*db.Queries, *sql.DB) {
	_, span := tracer.GetTracer().Start(ctx, "open")
	defer span.End()

	d, err := sql.Open("sqlite", "sqlitedb")
	if err != nil {
		log.Println(err)
		span.RecordError(err)
	}
	q := db.New(d)
	return q, d

}

func lookUPId(ctx context.Context, url string) (string, error) {
	octx, span := tracer.GetTracer().Start(ctx, "lookUpId")
	defer span.End()

	query, _ := open(octx)
	key, err := query.GetUrlId(octx, url)
	println(key)
	if err != nil {
		log.Println(err)
		span.RecordError(err)
		return "", err
	}
	sKey := strconv.FormatInt(key, 10)
	return sKey, nil

}

// ListAll method returns a json encoded blob with all the created URL's
func ListAll(ctx context.Context) (any, error) {
	lctx, span := tracer.GetTracer().Start(ctx, "ListAll")
	defer span.End()

	query, _ := open(lctx)
	key, err := query.GetUrls(lctx)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return key, nil

}

// CreateEntry adds a new entry in the database based on the api path given
func CreateEntry(c *gin.Context, url string) {
	cctx, span := tracer.GetTracer().Start(c.Request.Context(), "CreateEntry")
	defer span.End()

	v, err := lookUPId(cctx, url)
	if err != nil {
		log.Println(err)
		span.RecordError(err)
	}

	if !isCached(cctx, v) || !isLTS(cctx, url) {
		log.Println("already created")
		span.AddEvent("already cached")
	}

	query, _ := open(c.Request.Context())
	log.Println(url)
	_, err = query.CreateUrl(cctx, url)

	if err != nil {
		log.Println(err)
		span.RecordError(err)
	}

	err = memstorage.SetValue(v, url)
	if err != nil {
		log.Println(err)
		span.RecordError(err)
	}

	c.Done()
}

func isLTS(ctx context.Context, url string) bool {
	ictx, span := tracer.GetTracer().Start(ctx, "isLTS")
	defer span.End()

	query, _ := open(ictx)
	_, err := query.GetUrlId(ictx, url)
	return err == nil

}

// isCached checks against valkey if a given string is present in it's memory
func isCached(ctx context.Context, url string) bool {
	_, span := tracer.GetTracer().Start(ctx, "isCached")
	defer span.End()

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
	rctx, span := tracer.GetTracer().Start(c.Request.Context(), "CreateEntry")
	defer span.End()

	var url string
	q, _ := open(rctx)
	if !isCached(rctx, path) {
		i := toInt64(path)
		d, _ := q.GetUrlData(rctx, i)
		url = d.Name
		id := strconv.FormatInt(d.ID, 10)
		log.Println("path is not cached, performing sql query to get the dst addr")
		err := memstorage.SetValue(id, url)
		if err != nil {
			log.Fatalln("Unable to create entry in valkey ->", err)
			span.RecordError(err)
		}
		c.Redirect(http.StatusPermanentRedirect, url)
	} else {
		url, err := memstorage.GetKey(path)
		if err != nil {
			log.Println(err)
			span.RecordError(err)
		}
		log.Println("url has been cached")
		c.Redirect(http.StatusPermanentRedirect, url)
	}

}
