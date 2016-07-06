package interfaces

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Tinker-Ware/gh-service/domain"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
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
}

type WebServiceHandler struct {
	GHInteractor GHInteractor
	Sessions     *sessions.CookieStore
}

func (handler WebServiceHandler) Login(res http.ResponseWriter, req *http.Request) {

	url, state := handler.GHInteractor.GHLogin()

	fmt.Println("State login " + state)

	http.Redirect(res, req, url, http.StatusTemporaryRedirect)

}

func (handler WebServiceHandler) Callback(res http.ResponseWriter, req *http.Request) {

	incomingState := req.FormValue("state")
	code := req.FormValue("code")

	state := ""
	user, err := handler.GHInteractor.GHCallback(code, state, incomingState)
	if err != nil {
		return
	}

	usrB, _ := json.Marshal(user)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(usrB)

}

func (handler WebServiceHandler) Root(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(htmlIndex))
}

func (handler WebServiceHandler) GetCurrentUser(res http.ResponseWriter, req *http.Request) {
	session, err := handler.Sessions.Get(req, "user")
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}

	usr := session.Values["user"]

	userS := usr.(string)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(userS))

}

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

func (handler WebServiceHandler) ShowRepos(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	username := vars["username"]

	repos, err := handler.GHInteractor.ShowRepos(username)

	if err != nil {
		res.WriteHeader(http.StatusNotFound)
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
