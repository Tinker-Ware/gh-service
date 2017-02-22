package usecases

import "github.com/Tinker-Ware/gh-service/domain"

type GHInteractor struct {
	GithubRepository GithubRepository
}

type GithubRepository interface {
	GetOauthURL() (string, string)
	GetToken(code, givenState, incomingStates string) (*domain.User, error)
	SetToken(token string)
	GetAllRepos(username string) ([]domain.Repository, error)
	GetRepo(username, reponame string) (*domain.Repository, error)
	CreateRepo(username, reponame, org string, private bool) (*domain.Repository, error)
	GetKey(username string, id int) (*domain.Key, error)
	ShowKeys(username string) ([]domain.Key, error)
	CreateKey(username string, key *domain.Key) error
	CreateFile(file domain.File, author domain.Author, username, repoName string) error
	AddFiles(files []domain.File, author domain.Author, username, reponame string) error
	GetUser(username string) (*domain.User, error)
	AddDeployKey(username, reponame string, key *domain.Key) error
}

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func (interactor GHInteractor) GHLogin() (string, string) {
	url, oauthStateString := interactor.GithubRepository.GetOauthURL()
	return url, oauthStateString

}

func (interactor GHInteractor) GHCallback(code, state, incomingState string) (*domain.User, error) {

	usr, err := interactor.GithubRepository.GetToken(code, state, incomingState)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func (interactor GHInteractor) ShowUser(username string) (*domain.User, error) {

	usr, err := interactor.GithubRepository.GetUser(username)
	if err != nil {
		return nil, err
	}

	return usr, nil
}
