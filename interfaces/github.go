package interfaces

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Tinker-Ware/gh-service/domain"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	ghoauth "golang.org/x/oauth2/github"
)

// GithubRepository receives a github client and performs the necessary operations
// within an user account
type GithubRepository struct {
	client      *github.Client
	oauthConfig *oauth2.Config
}

// NewGithubRepository initializes the GithubRepository
func NewGithubRepository(clientID, clientSecret string, scopes []string) (*GithubRepository, error) {

	oauth2client := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		Endpoint:     ghoauth.Endpoint,
	}

	rp := &GithubRepository{
		oauthConfig: oauth2client,
	}

	return rp, nil
}

// README defines a readme to be written in a repository
var README = domain.File{
	Path:    "README.md",
	Content: []byte("# This is a README"),
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// GetOauthURL returns the OATH URL based on the configuration
func (repo GithubRepository) GetOauthURL() (string, string) {
	oauthStateString := randSeq(10)
	url := repo.oauthConfig.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	return url, oauthStateString
}

// GetToken returns a token from Github with the code received from the GET request
func (repo GithubRepository) GetToken(code, givenState, incomingStates string) (*domain.User, error) {

	token, err := repo.oauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		// TODO: Log with interactor Logger not yet implemented
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err.Error())
		return nil, err
	}

	oauthClient := repo.oauthConfig.Client(oauth2.NoContext, token)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get("")
	if err != nil {
		return nil, errors.New("Cannot retrieve User data")
	}

	usr := domain.User{
		Username:    *user.Login,
		AccessToken: token.AccessToken,
	}

	return &usr, nil
}

// GetUser retrieves the github user information
func (repo *GithubRepository) GetUser(username string) (*domain.User, error) {
	user, _, err := repo.client.Users.Get(username)

	if err != nil {
		return nil, err
	}

	usr := domain.User{
		ID:       *user.ID,
		Username: *user.Login,
	}

	return &usr, nil

}

// SetToken sets the token for the client
func (repo *GithubRepository) SetToken(token string) {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)
	repo.client = client

}

// GetAllRepos returns all the repos from an user account
func (repo GithubRepository) GetAllRepos(username string) ([]domain.Repository, error) {

	// list all repositories for the authenticated user
	ghrepos, _, err := repo.client.Repositories.List("", nil)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	repos := []domain.Repository{}

	for _, repo := range ghrepos {

		r := domain.Repository{
			Name:        repo.Name,
			FullName:    repo.FullName,
			Description: repo.Description,
			Private:     repo.Private,
			HTMLURL:     repo.HTMLURL,
			CloneURL:    repo.CloneURL,
			SSHURL:      repo.SSHURL,
		}
		repos = append(repos, r)
	}

	return repos, nil
}

// GetRepo gets the information from a single repo
func (repo GithubRepository) GetRepo(username, reponame string) (*domain.Repository, error) {

	rp, _, err := repo.client.Repositories.Get(username, reponame)
	if err != nil {
		return nil, err
	}

	r := &domain.Repository{
		Name:        rp.Name,
		FullName:    rp.FullName,
		Description: rp.Description,
		Private:     rp.Private,
		HTMLURL:     rp.HTMLURL,
		CloneURL:    rp.CloneURL,
		SSHURL:      rp.SSHURL,
	}

	return r, nil
}

// CreateRepo creates a repository in the github user account
func (repo GithubRepository) CreateRepo(username, reponame, org string, private bool) (*domain.Repository, error) {
	rp := &github.Repository{
		Name:    github.String(reponame),
		Private: github.Bool(private),
	}
	rp, _, err := repo.client.Repositories.Create(org, rp)

	if err != nil {
		return nil, err
	}

	r := &domain.Repository{
		Name:        rp.Name,
		FullName:    rp.FullName,
		Description: rp.Description,
		Private:     rp.Private,
		HTMLURL:     rp.HTMLURL,
		CloneURL:    rp.CloneURL,
		SSHURL:      rp.SSHURL,
	}

	return r, nil
}

// GetKey returns a Key from the user github account
func (repo GithubRepository) GetKey(username string, id int) (*domain.Key, error) {
	ghkey, _, err := repo.client.Users.GetKey(id)
	if err != nil {
		return nil, err
	}

	key := &domain.Key{
		ID:    ghkey.ID,
		Title: ghkey.Title,
		Key:   ghkey.Key,
		URL:   ghkey.URL,
	}

	return key, nil
}

// ShowKeys returns all the keys in a user github account
func (repo GithubRepository) ShowKeys(username string) ([]domain.Key, error) {
	ghKeys, _, err := repo.client.Users.ListKeys(username, nil)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	keys := []domain.Key{}

	for _, k := range ghKeys {
		key := domain.Key{
			ID:    k.ID,
			Key:   k.Key,
			Title: k.Title,
			URL:   k.URL,
		}
		keys = append(keys, key)
	}

	return keys, nil

}

// CreateKey creates a key in the github repository
func (repo GithubRepository) CreateKey(username string, key *domain.Key) error {

	k := github.Key{
		Title: key.Title,
		Key:   key.Key,
	}

	ghK, _, err := repo.client.Users.CreateKey(&k)
	if err != nil {
		return err
	}

	key.ID = ghK.ID
	key.URL = ghK.URL

	return nil

}

// CreateFile creates file inside an user repository
func (repo GithubRepository) CreateFile(file domain.File, author domain.Author, username, repoName string) error {

	opt := &github.RepositoryContentFileOptions{
		Message: github.String(author.Message),
		Content: file.Content,
		Branch:  github.String(author.Branch),
		Committer: &github.CommitAuthor{
			Name:  github.String(author.Author),
			Email: github.String(author.Email),
		},
	}

	_, _, err := repo.client.Repositories.CreateFile(username, repoName, file.Path, opt)
	if err != nil {
		return err
	}
	return nil
}

// AddFiles inserts multiple files inside a github repository
func (repo GithubRepository) AddFiles(files []domain.File, author domain.Author, username, reponame string) error {
	tree := []github.TreeEntry{}
	emptyRepo := "409 Git Repository is empty"
	lastCommit := ""

	// Get the reference sha
	// TODO remove hardcoded branch name
	ghTree, _, err := repo.client.Git.GetRef(username, reponame, "heads/master")
	if err != nil {

		// if the repo is empty create a README, else unexpected error
		if strings.Contains(err.Error(), emptyRepo) {

			// Create a copy of author for the initial commit
			author2 := author
			author2.Message = "Initial Commit"

			err = repo.CreateFile(README, author2, username, reponame)
			if err != nil {
				return err
			}

			// Repository should not be empty, otherwise, unexpected error
			ghTree, _, err = repo.client.Git.GetRef(username, reponame, "heads/master")
			if err != nil {
				return err
			}

		} else {
			return err

		}

	}
	lastCommit = *ghTree.Ref

	// Create a new tree
	// TODO Check for existing files and paths

	for _, file := range files {
		t := github.TreeEntry{
			Path:    github.String(file.Path),
			Mode:    github.String("100644"),
			Type:    github.String("blob"),
			Content: github.String(string(file.Content)),
		}

		tree = append(tree, t)
	}

	// Get the current tree
	currentTree, _, err := repo.client.Git.GetTree(username, reponame, lastCommit, false)
	if err != nil {
		return err
	}

	// Add the existing files to the tree
	for _, file := range currentTree.Entries {
		t := github.TreeEntry{
			Path: file.Path,
			Mode: file.Mode,
			Type: file.Type,
			SHA:  file.SHA,
		}

		tree = append(tree, t)
	}

	newTree, _, err := repo.client.Git.CreateTree(username, reponame, *currentTree.SHA, tree)
	if err != nil {
		return err
	}

	commit := github.Commit{
		Message: github.String(author.Message),
		Parents: []github.Commit{{SHA: currentTree.SHA}},
		Tree:    &github.Tree{SHA: newTree.SHA},
	}

	newCommit, _, err := repo.client.Git.CreateCommit(username, reponame, &commit)
	if err != nil {
		return err
	}

	reference := github.Reference{
		Ref: github.String("refs/heads/master"),
		Object: &github.GitObject{
			SHA: newCommit.SHA,
		},
	}

	_, _, err = repo.client.Git.UpdateRef(username, reponame, &reference, false)
	if err != nil {
		return err
	}

	return nil
}

func randSeq(n int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
