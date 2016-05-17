package interfaces

import (
	"log"
	"net/http"
	"time"

	"github.com/Tinker-Ware/gh-service/domain"
)

type Adapter func(http.Handler) http.Handler
type repository interface {
	SetToken(token string)
}

func Notify() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			defer log.Printf("%s on %s took %s\n", r.Method, r.URL.Path, time.Since(start))
			h.ServeHTTP(w, r)
		})
	}
}

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func SetToken(repo repository) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get(domain.TokenHeader)
			repo.SetToken(token)
			h.ServeHTTP(w, r)
		})
	}
}
