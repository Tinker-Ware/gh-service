package interfaces

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Tinker-Ware/gh-service/domain"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

const htmlIndex = `<html><body>
Logged in with <a href="/github/login">GitHub</a>
</body></html>`

const htmlCloseWindow = `<html><body>
Logged in with <a href="/login">GitHub</a>
</body></html>`

type repoRequest struct {
	Owner   string `json:"owner"`
	Name    string `json:"name"`
	Private bool   `json:"private"`
	Org     string `json:"org"`
}

type fileRequest struct {
	domain.Author
	domain.File
}

type multipleFilesRequest struct {
	Author domain.Author `json:"author"`
	Files  []domain.File `json:"files"`
}

type repositoryResponse struct {
	Repository *domain.Repository `json:"repository"`
}

type repositoriesResponse struct {
	Repositories []domain.Repository `json:"repositories"`
}

// GHInteractor defines all the functions the Github Interactor should have
type GHInteractor interface {
	GHCallback(code, state, incomingState string) (*domain.User, error)
	GHLogin() (string, string)
	ShowUser(username string) (*domain.User, error)
	ShowRepos(username string) ([]domain.Repository, error)
	CreateRepo(username, reponame, org string, private bool) (*domain.Repository, error)
	ShowRepo(username, repo string) (*domain.Repository, error)
	ShowKeys(username string) ([]domain.Key, error)
	CreateKey(username string, key *domain.Key) error
	ShowKey(username string, id int) (*domain.Key, error)
	CreateFile(file domain.File, author domain.Author, username, repo string) error
	AddFiles(files []domain.File, author domain.Author, username, repo string) error
	AddDeployKey(username, reponame string, key *domain.Key) error
}

// WebServiceHandler has all the necessary fields to run a web-based interface
type WebServiceHandler struct {
	GHInteractor GHInteractor
	APIHost      string
}

// Login is a helper method to test the Github oauth login
func (handler WebServiceHandler) Login(res http.ResponseWriter, req *http.Request) {

	url, state := handler.GHInteractor.GHLogin()

	fmt.Println("State login " + state)

	http.Redirect(res, req, url, http.StatusTemporaryRedirect)

}

type oauthWrapper struct {
	OauthRequest oauthRequest `json:"oauth_request"`
}

type oauthRequest struct {
	UserID int    `json:"user_id"`
	Code   string `json:"code"`
	State  string `json:"state"`
}

type integrationWrapper struct {
	Integration integration `json:"integration"`
}

type callbackResponse struct {
	Callback callback `json:"callback"`
}

type callback struct {
	Provider string `json:"provider"`
	Username string `json:"username"`
}

type httpError struct {
	Error string `json:"error"`
}

const integrationURL string = "/api/v1/users/%d/integration"

// Callback manages the Github OAUTH callback
// TODO: refactor this function, it has grown to big
func (handler WebServiceHandler) Callback(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	userToken := req.Header.Get("Authorization")

	decoder := json.NewDecoder(req.Body)

	var oauthwrapper oauthWrapper

	err := decoder.Decode(&oauthwrapper)
	if err != nil {
		res.WriteHeader(422)

		errS := fmt.Sprintf("cannot process request %s", err.Error())

		log.Println(errS)

		resErr := httpError{
			Error: "cannot process request",
		}

		respBytes, _ := json.Marshal(resErr)
		res.Header().Set("Content-Type", "application/json")
		res.Write(respBytes)
		return
	}

	token, err := handler.GHInteractor.GHCallback(oauthwrapper.OauthRequest.Code, "", oauthwrapper.OauthRequest.State)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		errS := fmt.Sprintf("Github oauth error: %s", err.Error())

		log.Println(errS)

		resErr := httpError{
			Error: errS,
		}

		respBytes, _ := json.Marshal(resErr)
		res.Header().Set("Content-Type", "application/json")
		res.Write(respBytes)
		return
	}

	wrapper := integrationWrapper{
		Integration: integration{
			UserID:     oauthwrapper.OauthRequest.UserID,
			Token:      token.AccessToken,
			Provider:   "github",
			Username:   token.Username,
			ExpireDate: token.ExpirationDate,
		},
	}

	reqBytes, _ := json.Marshal(&wrapper)

	buf := bytes.NewBuffer(reqBytes)

	path := fmt.Sprintf(integrationURL, oauthwrapper.OauthRequest.UserID)

	request, _ := http.NewRequest(http.MethodPost, handler.APIHost+path, buf)
	request.Header.Add("Authorization", userToken)
	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	resp, _ := client.Do(request)
	if resp.StatusCode != http.StatusCreated {
		res.WriteHeader(http.StatusInternalServerError)
		resErr := httpError{
			Error: "cannot save integration",
		}

		respBytes, _ := json.Marshal(resErr)
		res.Header().Set("Content-Type", "application/json")
		res.Write(respBytes)
		return
	}

	response := callbackResponse{
		Callback: callback{
			Provider: "github",
			Username: token.Username,
		},
	}

	respBytes, _ := json.Marshal(&response)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(respBytes)

}

// Root is a function for the index of the microservice
func (handler WebServiceHandler) Root(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(htmlIndex))
}

// ShowUser returns the current logged user
func (handler WebServiceHandler) ShowUser(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	username := vars["username"]

	user, err := handler.GHInteractor.ShowUser(username)

	if err != nil {
		fmt.Println(err.Error())
		res.WriteHeader(http.StatusNotFound)
		return
	}

	userB, _ := json.Marshal(user)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(userB))
}

// ShowRepos returns a JSON with all the repositories within an user account
func (handler WebServiceHandler) ShowRepos(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	username := vars["username"]

	repos, err := handler.GHInteractor.ShowRepos(username)

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		errS := fmt.Sprintf("Cannot retrieve repositories: %s", err.Error())

		log.Println(errS)

		resErr := httpError{
			Error: errS,
		}

		respBytes, _ := json.Marshal(resErr)
		res.Header().Set("Content-Type", "application/json")
		res.Write(respBytes)
		return
	}

	reposR := repositoriesResponse{
		Repositories: repos,
	}

	reposB, _ := json.Marshal(reposR)

	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(reposB))
}

// CreateRepo creates a repository in the user account and returns a JSON response
func (handler WebServiceHandler) CreateRepo(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	decoder := json.NewDecoder(req.Body)
	repo := repoRequest{}
	err := decoder.Decode(&repo)
	if err != nil {
		fmt.Println(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	r, err := handler.GHInteractor.CreateRepo(repo.Owner, repo.Name, repo.Org, repo.Private)
	if err != nil {
		fmt.Println(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	repoB, _ := json.Marshal(r)

	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(repoB))

}

// ShowRepo returns a JSON response of a single repository
func (handler WebServiceHandler) ShowRepo(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	username := vars["username"]
	repoName := vars["repo"]

	repo, err := handler.GHInteractor.ShowRepo(username, repoName)
	if err != nil {
		fmt.Println(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return

	}

	repoR := repositoryResponse{
		Repository: repo,
	}

	repoB, _ := json.Marshal(repoR)

	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(repoB))

}

// ShowKeys returns all the keys within github in a JSON response
func (handler WebServiceHandler) ShowKeys(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	username := vars["username"]

	keys, err := handler.GHInteractor.ShowKeys(username)
	if err != nil {
		fmt.Println(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return

	}

	keysB, _ := json.Marshal(keys)

	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(keysB))

}

// CreateKey creates a key in the user repository
func (handler WebServiceHandler) CreateKey(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	decoder := json.NewDecoder(req.Body)
	key := domain.Key{}
	err := decoder.Decode(&key)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(req)
	username := vars["username"]

	err = handler.GHInteractor.CreateKey(username, &key)

	if err != nil {
		fmt.Println(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	keyB, _ := json.Marshal(key)

	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(keyB))
}

// ShowKey shows a single key whithin the user account on github
func (handler WebServiceHandler) ShowKey(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	username := vars["username"]
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	key, err := handler.GHInteractor.ShowKey(username, id)
	if err != nil {
		fmt.Println(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	keyB, _ := json.Marshal(key)

	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(keyB))

}

// AddFileToRepository adds a single file within an user repository
func (handler WebServiceHandler) AddFileToRepository(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	username := vars["username"]
	repoName := vars["repo"]

	decoder := json.NewDecoder(req.Body)
	file := fileRequest{}
	err := decoder.Decode(&file)
	if err != nil {
		fmt.Println(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = handler.GHInteractor.CreateFile(file.File, file.Author, username, repoName)
	if err != nil {
		fmt.Println(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusCreated)

}

// AddMultipleFilesToRepository add multiple files into a repository
func (handler WebServiceHandler) AddMultipleFilesToRepository(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	vars := mux.Vars(req)
	username := vars["username"]
	repoName := vars["repo"]

	decoder := json.NewDecoder(req.Body)

	request := multipleFilesRequest{}

	err := decoder.Decode(&request)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = handler.GHInteractor.AddFiles(request.Files, request.Author, username, repoName)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusCreated)

}

type keyWrapper struct {
	Key domain.Key `json:"deploy_key"`
}

func (handler WebServiceHandler) CreateRepoDeployKey(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	vars := mux.Vars(req)

	defer req.Body.Close()

	username := vars["username"]
	repoName := vars["repo"]

	decoder := json.NewDecoder(req.Body)
	var key keyWrapper
	err := decoder.Decode(&key)
	if err != nil {
		log.Printf("Error unmarshaling json in CreateRepoDeployKey %s", err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = handler.GHInteractor.AddDeployKey(username, repoName, &key.Key)
	if err != nil {
		log.Printf("Cannot create deploy key %s", err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	keyB, _ := json.Marshal(key)

	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(keyB))

}

func tokenToJSON(token *oauth2.Token) (string, error) {
	d, err := json.Marshal(token)
	if err != nil {
		return "", err
	}
	return string(d), nil

}

func tokenFromJSON(jsonStr string) (*oauth2.Token, error) {
	var token oauth2.Token
	if err := json.Unmarshal([]byte(jsonStr), &token); err != nil {
		return nil, err
	}
	return &token, nil
}
