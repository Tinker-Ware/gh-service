package usecases

import "github.com/Tinker-Ware/gh-service/domain"

func (interactor GHInteractor) ShowKeys(username string) ([]domain.Key, error) {

	keys, err := interactor.GithubRepository.ShowKeys(username)
	if err != nil {
		return nil, err
	}

	return keys, nil

}

func (interactor GHInteractor) CreateKey(username string, key *domain.Key) error {

	err := interactor.GithubRepository.CreateKey(username, key)
	if err != nil {
		return err
	}
	return nil
}

func (interactor GHInteractor) ShowKey(username string, id int) (*domain.Key, error) {

	key, err := interactor.GithubRepository.GetKey(username, id)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (interactor GHInteractor) AddDeployKey(username, reponame string, key *domain.Key) error {
	err := interactor.GithubRepository.AddDeployKey(username, reponame, key)
	return err
}
