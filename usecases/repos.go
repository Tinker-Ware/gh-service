package usecases

import "github.com/gh-service/domain"

func (interactor GHInteractor) ShowRepos(username string) ([]domain.Repository, error) {

	repos, err := interactor.GithubRepository.GetAllRepos(username)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

func (interactor GHInteractor) CreateRepo(username, token, reponame, org string, private bool) (*domain.Repository, error) {
	r, err := interactor.GithubRepository.CreateRepo(username, reponame, org, private)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (interactor GHInteractor) ShowRepo(username, token, repo string) (*domain.Repository, error) {

	r, err := interactor.GithubRepository.GetRepo(username, repo)
	if err != nil {
		return nil, err
	}
	return r, nil
}
