package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"

	"github.com/Tinker-Ware/gh-service/infrastructure"
	"github.com/Tinker-Ware/gh-service/interfaces"
	"github.com/Tinker-Ware/gh-service/usecases"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const defaultPath = "/etc/gh-service.conf"

// Define configuration flags
var confFilePath = flag.String("conf", defaultPath, "Custom path for configuration file")

func main() {

	flag.Parse()

	config, err := infrastructure.GetConfiguration(*confFilePath)
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot parse configuration")
	}

	ghrepo, err := interfaces.NewGithubRepository(config.ClientID, config.ClientSecret, config.Scopes)
	if err != nil {
		panic(err.Error())
	}

	ghinteractor := usecases.GHInteractor{
		GithubRepository: ghrepo,
	}

	handler := interfaces.WebServiceHandler{
		GHInteractor: ghinteractor,
		APIHost:      config.APIHost,
	}

	// Add CORS Support
	headers := handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Authorization"})
	origins := handlers.AllowedOrigins([]string{"http://localhost", "http://provision.tinkerware.io", "https://provision.tinkerware.io"})

	r := mux.NewRouter()

	subrouter := r.PathPrefix("/api/v1/repository/github").Subrouter()
	subrouter.Handle("/oauth", interfaces.Adapt(http.HandlerFunc(handler.Callback), interfaces.Notify())).Methods("POST")
	subrouter.Handle("/{username}/repos", interfaces.Adapt(http.HandlerFunc(handler.ShowRepos), interfaces.Notify(), interfaces.GetToken(ghrepo, config.APIHost, config.Salt))) //.Methods("GET")
	subrouter.Handle("/{username}/{repo}", interfaces.Adapt(http.HandlerFunc(handler.ShowRepo), interfaces.Notify(), interfaces.GetToken(ghrepo, config.APIHost, config.Salt))) //.Methods("GET")
	subrouter.Handle("/{username}/{repo}/deploy_key", interfaces.Adapt(http.HandlerFunc(handler.CreateRepoDeployKey), interfaces.Notify(), interfaces.GetToken(ghrepo, config.APIHost, config.Salt))).Methods("POST")
	// subrouter.Handle("/user/{username}/repos", interfaces.Adapt(http.HandlerFunc(handler.CreateRepo), interfaces.Notify(), interfaces.SetToken(ghrepo))).Methods("POST")
	// subrouter.Handle("/user/{username}/keys", interfaces.Adapt(http.HandlerFunc(handler.CreateRepo), interfaces.Notify(), interfaces.SetToken(ghrepo))).Methods("GET")
	// subrouter.Handle("/user/{username}/keys", interfaces.Adapt(http.HandlerFunc(handler.CreateKey), interfaces.Notify(), interfaces.SetToken(ghrepo))).Methods("POST")
	// subrouter.Handle("/user/{username}/keys/{id}", interfaces.Adapt(http.HandlerFunc(handler.ShowKey), interfaces.Notify(), interfaces.SetToken(ghrepo))).Methods("GET")
	// subrouter.Handle("/user/{username}/{repo}/addfile", interfaces.Adapt(http.HandlerFunc(handler.AddFileToRepository), interfaces.Notify(), interfaces.SetToken(ghrepo))).Methods("POST")
	// subrouter.Handle("/user/{username}/{repo}/addfiles", interfaces.Adapt(http.HandlerFunc(handler.AddMultipleFilesToRepository), interfaces.Notify(), interfaces.SetToken(ghrepo))).Methods("POST")
	// subrouter.HandleFunc("/user_info", handler.GetCurrentUser).Methods("GET")

	n := negroni.Classic()
	n.UseHandler(handlers.CORS(headers, origins)(r))

	port := bytes.Buffer{}

	port.WriteString(":")
	port.WriteString(config.Port)

	n.Run(port.String())
}
