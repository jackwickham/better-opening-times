package main

import (
	"fmt"
	"net/http"
	"time"
)

type CacheMiddleware struct {
	duration time.Duration
	next     http.Handler
}

func (m CacheMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", fmt.Sprintf("max-age=%d, public", int64(m.duration.Seconds())))
	m.next.ServeHTTP(w, r)
}

func WrapCache(duration time.Duration, handler http.Handler) CacheMiddleware {
	return CacheMiddleware{
		duration: duration,
		next:     handler,
	}
}
