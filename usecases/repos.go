package usecases

import "github.com/Tinker-Ware/gh-service/domain"

func (interactor GHInteractor) ShowRepos(username string) ([]domain.Repository, error) {

	repos, err := interactor.GithubRepository.GetAllRepos(username)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

func (interactor GHInteractor) CreateRepo(username, reponame, org string, private bool) (*domain.Repository, error) {
	r, err := interactor.GithubRepository.CreateRepo(username, reponame, org, private)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (interactor GHInteractor) ShowRepo(username, repo string) (*domain.Repository, error) {

	r, err := interactor.GithubRepository.GetRepo(username, repo)
	if err != nil {
		return nil, err
	}
	return r, nil
}
