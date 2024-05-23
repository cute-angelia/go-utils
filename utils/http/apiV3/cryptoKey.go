package apiV3

import (
	"context"
	"net/http"
)

// CryptoEr 从中间件设置加密密钥
var CryptoEr = cryptoEr{}

type cryptoEr struct {
}

var (
	CryptoCtxKey = &contextKey{"CryptoKey"}
)

func (that cryptoEr) SetCryptoKey(cryptoKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), CryptoCtxKey, cryptoKey))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// GetRequestContentType is a helper function that returns ContentType based on
// context or request headers.
func (that cryptoEr) GetRequestContentType(r *http.Request) string {
	if value, ok := r.Context().Value(CryptoCtxKey).(string); ok {
		return value
	}
	return ""
}
