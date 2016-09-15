package interfaces

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Tinker-Ware/gh-service/domain"
	"github.com/dvsekhvalnov/jose2go"
)

// Adapter is the signature of an HTTPHandler for middlewares
type Adapter func(http.Handler) http.Handler
type repository interface {
	SetToken(token string)
}

// Notify is a middleware to measure the time that a request takes
func Notify() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			defer log.Printf("%s on %s took %s\n", r.Method, r.URL.Path, time.Since(start))
			h.ServeHTTP(w, r)
		})
	}
}

// Adapt takes several Adapters and calls them in order
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

// SetToken injects the token from the request
func SetToken(repo repository) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get(domain.TokenHeader)
			repo.SetToken(token)
			h.ServeHTTP(w, r)
		})
	}
}

type integration struct {
	UserID     int    `json:"user_id"`
	Token      string `json:"token"`
	Username   string `json:"username"`
	Provider   string `json:"provider"`
	ExpireDate int64  `json:"expire_date"`
}

type integrationsWrapper struct {
	Integrations []integration `json:"integrations"`
}

const integrationsURL string = "/api/v1/users/%s/integrations"

// GetToken gets the token from the users microservice
func GetToken(repo repository, apiURL string, salt string) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userToken := r.Header.Get("authorization")

			payload, _, _ := jose.Decode(userToken, []byte(salt))

			var objmap map[string]*json.RawMessage
			json.Unmarshal([]byte(payload), &objmap)

			var aud string
			json.Unmarshal(*objmap["aud"], &aud)

			vals := strings.Split(aud, ":")

			path := fmt.Sprintf(integrationsURL, vals[1])

			request, _ := http.NewRequest(http.MethodGet, apiURL+path, nil)
			request.Header.Set("authorization", userToken)

			client := &http.Client{}
			resp, _ := client.Do(request)
			defer resp.Body.Close()

			integrations := integrationsWrapper{}
			decoder := json.NewDecoder(resp.Body)
			decoder.Decode(&integrations)
			var token string

			for _, integ := range integrations.Integrations {
				if integ.Provider == "github" {
					token = integ.Token
				}
			}

			repo.SetToken(token)

			h.ServeHTTP(w, r)
		})
	}
}
