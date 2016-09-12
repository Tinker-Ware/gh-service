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

type integration struct {
	Token    string `json:"token"`
	Provider string `json:"provider"`
}

type integrationWrapper struct {
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

			integrations := integrationWrapper{}
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
