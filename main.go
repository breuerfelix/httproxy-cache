package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	buffer *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.buffer.Write(b)
	return w.ResponseWriter.Write(b)
}

type Response struct {
	Body   []byte
	Header http.Header
	Status int
}

func main() {
	config := NewConfig()
	cache := NewInMemoryCache[Response]()

	target, err := url.Parse(config.Upstream)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	proxy := httputil.NewSingleHostReverseProxy(target)

	for _, route := range config.Routes {
		router.Any(route.Path, func(c *gin.Context) {
			key := c.Request.URL.Path

			if c.Request.Method == http.MethodHead {
				log.Println("clearing cache with HEAD")
				cache.Delete(key)
				c.Status(http.StatusNoContent)
				return
			}

			for _, method := range route.ClearCache {
				if strings.EqualFold(c.Request.Method, method) {
					log.Println("clearing cache...")
					cache.Delete(key)
					break
				}
			}

			if c.Request.Method != http.MethodGet {
				proxy.ServeHTTP(c.Writer, c.Request)
				return
			}

			cached := cache.Get(key)
			if cached != nil {
				// copy headers
				dst := c.Writer.Header()
				for k, vv := range cached.Header {
					for _, v := range vv {
						dst.Add(k, v)
					}
				}
				c.Writer.WriteHeader(cached.Status)
				c.Writer.Write(cached.Body)
				log.Println("sent from cache")
				return
			}

			buffer := &bytes.Buffer{}
			c.Writer = responseBodyWriter{c.Writer, buffer}

			log.Println("sent from upstream")
			proxy.ServeHTTP(c.Writer, c.Request)

			log.Println("saved to cache")
			go cache.Set(key, &Response{
				Status: c.Writer.Status(),
				Header: c.Writer.Header().Clone(),
				Body:   buffer.Bytes(),
			})
		})
	}

	// proxy all other requests
	router.NoRoute(func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	router.Run(":8080")
}
