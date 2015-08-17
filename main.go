package main

import (
	"flag"
	"fmt"

	"golang.org/x/oauth2"

	"github.com/codegangsta/negroni"
	"github.com/gh-service/infraestructure"
	"github.com/gh-service/interfaces"
	"github.com/gh-service/usecases"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	ghoauth "golang.org/x/oauth2/github"
)

// TODO: DRY client usage in handlers
// TODO: Cache to avoid multiple calls to the GH API
// TODO: Pass the token in the header to delegate the management to another microservice
// TODO: Once the token is passed via header change the routes to add org support in repo listing
// TODO: Remove struct types in Webservice and use the domain structs
// TODO: Figure out how to get a token without OAUTH to use tests
// TODO: Inject GH API data from here

// Define configuration flags
var confFilePath = flag.String("conf", "/etc/iaas/gh-service.yaml", "Custom path for configuration file")

func main() {

	flag.Parse()

	config, err := infraestructure.GetConfiguration(*confFilePath)
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot parse configuration")
	}

	oauth2client := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Scopes:       config.Scopes,
		Endpoint:     ghoauth.Endpoint,
	}

	userRepo := interfaces.UserRepo{}
	ghinteractor := usecases.GHInteractor{
		UserRepo:    userRepo,
		OauthConfig: oauth2client,
	}

	store := sessions.NewCookieStore([]byte("something-very-secret"))

	handler := interfaces.WebServiceHandler{
		GHInteractor: ghinteractor,
		Sessions:     store,
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handler.Root)
	r.HandleFunc("/login", handler.Login)
	r.HandleFunc("/github_oauth_cb", handler.Callback)
	r.HandleFunc("/user/{username}", handler.ShowUser).Methods("GET")
	r.HandleFunc("/user/{username}/repos", handler.ShowRepos).Methods("GET")
	r.HandleFunc("/user/{username}/repos", handler.CreateRepo).Methods("POST")
	r.HandleFunc("/user/{username}/keys", handler.ShowKeys).Methods("GET")
	r.HandleFunc("/user/{username}/keys", handler.CreateKey).Methods("POST")
	r.HandleFunc("/user/{username}/keys/{id}", handler.ShowKey).Methods("GET")
	r.HandleFunc("/user/{username}/{repo}", handler.ShowRepo).Methods("GET")
	r.HandleFunc("/user_info", handler.GetCurrentUser).Methods("GET")

	n := negroni.Classic()
	n.UseHandler(r)
	n.Run(":" + config.Port)
}
