package apiV3

import (
	"bytes"
	"github.com/go-chi/stampede"
	"io"
	"net/http"
	"strings"
	"time"
)

var Stampeder = stampeder{}

type stampeder struct {
}

func (that stampeder) CachedPost(cacheSize int, ttl time.Duration) func(next http.Handler) http.Handler {
	return stampede.HandlerWithKey(cacheSize, ttl, func(r *http.Request) uint64 {
		// Read the request payload, and then setup buffer for future reader
		var buf []byte
		if r.Body != nil {
			buf, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(buf))
		}
		token := r.Header.Get("Authorization")
		return stampede.BytesToHash(
			[]byte(strings.ToLower(r.URL.Path)),
			[]byte(strings.ToLower(token)),
			buf,
		)
	})
}
