package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"url_shortener/router"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	r := router.SetRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthz", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestIndex(t *testing.T) {
	r := router.SetRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestAddRoute(t *testing.T) {
	r := router.SetRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/addRoute", strings.NewReader(`{"url": "https://example.com"}`))
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestListRoute(t *testing.T) {
	r := router.SetRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/list", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
