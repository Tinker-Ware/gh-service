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
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
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

	store := sessions.NewCookieStore([]byte("something-very-secret"))

	handler := interfaces.WebServiceHandler{
		GHInteractor: ghinteractor,
		Sessions:     store,
	}

	r := mux.NewRouter()

	subrouter := r.PathPrefix("/api/v1/repository/github").Subrouter()
	subrouter.Handle("/{username}/repos", interfaces.Adapt(http.HandlerFunc(handler.ShowRepos), interfaces.Notify(), interfaces.GetToken(ghrepo, config.Host, config.Salt))).Methods("GET")
	subrouter.Handle("/{username}/{repo}", interfaces.Adapt(http.HandlerFunc(handler.ShowRepo), interfaces.Notify(), interfaces.GetToken(ghrepo, config.Host, config.Salt))).Methods("GET")
	// subrouter.Handle("/user/{username}", interfaces.Adapt(http.HandlerFunc(handler.ShowUser), interfaces.Notify(), interfaces.SetToken(ghrepo))).Methods("GET")
	// subrouter.Handle("/user/{username}/repos", interfaces.Adapt(http.HandlerFunc(handler.CreateRepo), interfaces.Notify(), interfaces.SetToken(ghrepo))).Methods("POST")
	// subrouter.Handle("/user/{username}/keys", interfaces.Adapt(http.HandlerFunc(handler.CreateRepo), interfaces.Notify(), interfaces.SetToken(ghrepo))).Methods("GET")
	// subrouter.Handle("/user/{username}/keys", interfaces.Adapt(http.HandlerFunc(handler.CreateKey), interfaces.Notify(), interfaces.SetToken(ghrepo))).Methods("POST")
	// subrouter.Handle("/user/{username}/keys/{id}", interfaces.Adapt(http.HandlerFunc(handler.ShowKey), interfaces.Notify(), interfaces.SetToken(ghrepo))).Methods("GET")
	// subrouter.Handle("/user/{username}/{repo}/addfile", interfaces.Adapt(http.HandlerFunc(handler.AddFileToRepository), interfaces.Notify(), interfaces.SetToken(ghrepo))).Methods("POST")
	// subrouter.Handle("/user/{username}/{repo}/addfiles", interfaces.Adapt(http.HandlerFunc(handler.AddMultipleFilesToRepository), interfaces.Notify(), interfaces.SetToken(ghrepo))).Methods("POST")
	// subrouter.HandleFunc("/user_info", handler.GetCurrentUser).Methods("GET")

	n := negroni.Classic()
	n.UseHandler(r)

	port := bytes.Buffer{}

	port.WriteString(":")
	port.WriteString(config.Port)

	n.Run(port.String())
}
